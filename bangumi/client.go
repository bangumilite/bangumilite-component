package bangumi

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/sstp105/bangumi-component/httplib"
	"github.com/sstp105/bangumi-component/model"
	"net/http"
	"strings"
	"sync"
)

type APIPath string
type APIError string

const (
	HTMLBaseURL string = "https://bangumi.tv"                // For scraping HTML content
	APIBaseURL  string = "https://api.bgm.tv"                // For bangumi RESTful APIs
	OAuthURL    string = "https://bgm.tv/oauth/access_token" // For bangumi OAuth flow

	UserAgentHeader           = "sstp105/bangumi-services (GCP; Golang; Private Project)"
	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"

	MaxConcurrentGoroutines = 10 // Max goroutines for concurrent API calls

	APIPathGetSubject           APIPath = "/v0/subjects/%d"
	APIPathGetSubjectCharacters APIPath = "/v0/subjects/%d/characters"

	ErrorGeneric APIError = "ErrorGeneric"
	ErrorOAuth   APIError = "ErrorOAuth"
)

type Client struct {
	client *resty.Client
	logger *logrus.Logger
}

func NewClient(logger *logrus.Logger) *Client {
	client := httplib.NewClient()

	return &Client{
		client: client,
		logger: logger,
	}
}

func (c *Client) GetSubjects(ctx context.Context, ids []int) ([]model.BangumiSubject, error) {
	return c.GetSubjectsWithAccessToken(ctx, ids, "")
}

func (c *Client) GetSubjectsWithAccessToken(ctx context.Context, ids []int, accessToken string) ([]model.BangumiSubject, error) {
	return concurrentFetch(ctx, ids, func(ctx context.Context, id int) (model.BangumiSubject, error) {
		subject, err := c.GetSubjectWithAccessToken(ctx, id, accessToken)
		if err != nil {
			return model.BangumiSubject{}, err
		}

		return *subject, nil
	})
}

func (c *Client) GetSubject(ctx context.Context, id int) (*model.BangumiSubject, error) {
	return c.GetSubjectWithAccessToken(ctx, id, "")
}

func (c *Client) GetSubjectWithAccessToken(ctx context.Context, id int, accessToken string) (*model.BangumiSubject, error) {
	url := apiURL(APIPathGetSubject, id)
	subject := model.BangumiSubject{}

	req := c.client.R().
		SetContext(ctx).
		SetHeader("User-Agent", UserAgentHeader).
		SetHeader("Content-Type", ContentTypeJSON).
		SetResult(&subject).
		SetError(model.BangumiGenericErrorResponse{})

	if accessToken != "" {
		req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, handleError(resp, url, ErrorGeneric)
	}

	return &subject, nil
}

func (c *Client) GetSubjectCharacters(ctx context.Context, id int) ([]model.BangumiRelatedCharacter, error) {
	url := apiURL(APIPathGetSubjectCharacters, id)
	var characters []model.BangumiRelatedCharacter

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("User-Agent", UserAgentHeader).
		SetHeader("Content-Type", ContentTypeJSON).
		SetResult(&characters).
		SetError(model.BangumiGenericErrorResponse{}).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, handleError(resp, url, ErrorGeneric)
	}

	return characters, nil
}

func (c *Client) RefreshAccessToken(ctx context.Context, token model.FirestoreBangumiToken) (*model.BangumiOAuthResponse, error) {
	tokenResp := model.BangumiOAuthResponse{}
	formData := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     token.ClientID,
		"client_secret": token.ClientSecret,
		"redirect_uri":  token.RedirectURI,
		"refresh_token": token.RefreshToken,
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", ContentTypeFormURLEncoded).
		SetFormData(formData).
		SetResult(&tokenResp).
		SetError(model.BangumiOAuthErrorResponse{}).
		Post(OAuthURL)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, handleError(resp, OAuthURL, ErrorOAuth)
	}

	return &tokenResp, nil
}

func (c *Client) GetHTML(ctx context.Context, path string) (*goquery.Document, error) {
	url := fmt.Sprintf("%s/%s", HTMLBaseURL, path)

	resp, err := c.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request failed, status code:%d", resp.StatusCode())
	}

	body := resp.String()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func apiURL(p APIPath, args ...interface{}) string {
	return fmt.Sprintf(APIBaseURL+string(p), args...)
}

// handleError is a generic function to handle errors from bangumi API responses.
func handleError(resp *resty.Response, path string, errorType APIError) error {
	switch errorType {
	case ErrorOAuth:
		errResp := resp.Error().(*model.BangumiOAuthErrorResponse)
		return apiError(path, resp.StatusCode(), fmt.Sprintf("error:%s,details:%s", errResp.Error, errResp.ErrorDescription))
	case ErrorGeneric:
		errResp := resp.Error().(*model.BangumiGenericErrorResponse)
		return apiError(path, resp.StatusCode(), fmt.Sprintf("error:%s,details:%s", errResp.Title, errResp.Description))
	default:
		return fmt.Errorf("unexpected error type: %s", errorType)
	}
}

func apiError(path string, statusCode int, message string) error {
	return fmt.Errorf("failed to call %s, status code: %d, error: %s",
		path,
		statusCode,
		message,
	)
}

// concurrentFetch fetches data concurrently for a list of IDs using a provided fetch function.
// It ensures that the fetches are limited by a maximum number of concurrent goroutines.
// The function will wait for all fetches to complete and return the results.
func concurrentFetch[T any](
	ctx context.Context,
	ids []int,
	fetchFunc func(ctx context.Context, id int) (T, error),
) ([]T, error) {
	var (
		results []T
		mu      sync.Mutex
		wg      sync.WaitGroup
		sem     = make(chan struct{}, MaxConcurrentGoroutines)
		errCh   = make(chan error, len(ids))
	)

	for _, id := range ids {
		wg.Add(1)
		sem <- struct{}{}

		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }()

			result, err := fetchFunc(ctx, id)
			if err != nil {
				errCh <- err
				return
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(id)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return nil, <-errCh
	}

	return results, nil
}

func GetVoiceActorsFromCharacters(characters []model.BangumiRelatedCharacter) []model.BangumiPerson {
	mp := make(map[int]bool)
	var actors []model.BangumiPerson

	for _, character := range characters {
		for _, actor := range character.Actors {
			if _, exists := mp[actor.ID]; !exists {
				mp[actor.ID] = true
				actors = append(actors, actor)
			}
		}
	}

	if len(actors) < 5 {
		return actors
	}

	return actors[:5]
}

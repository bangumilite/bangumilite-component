package bangumi

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bangumilite/bangumilite-component/httplib"
	"github.com/bangumilite/bangumilite-component/model"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"sync"
)

type APIPath string
type APIError string
type GrantType string
type RequestOption func(req *resty.Request)

const (
	HTMLBaseURL string = "https://bangumi.tv"
	APIBaseURL  string = "https://api.bgm.tv"
	OAuthURL    string = "https://bgm.tv/oauth/access_token"

	RefreshToken GrantType = "refresh_token"

	MonoPath string = "/mono"

	UserAgentHeader           = "github.com/bangumilite (GCP; Golang; Private Project)"
	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"

	MaxConcurrentGoroutines = 10

	APIPathGetSubject           APIPath = "/v0/subjects/%d"
	APIPathGetSubjectCharacters APIPath = "/v0/subjects/%d/characters"

	ErrorGeneric APIError = "ErrorGeneric"
	ErrorOAuth   APIError = "ErrorOAuth"
)

type Client struct {
	client *resty.Client
}

func NewClient() *Client {
	client := httplib.NewClient()
	return &Client{
		client: client,
	}
}

func (c *Client) GetSubjects(ctx context.Context, ids []int, opts ...RequestOption) ([]model.BangumiSubject, error) {
	return batchFetch(ctx, ids, func(ctx context.Context, id int) (model.BangumiSubject, error) {
		subject, err := c.GetSubject(ctx, id, opts...)
		if err != nil {
			return model.BangumiSubject{}, err
		}

		return *subject, nil
	})
}

func (c *Client) GetSubject(ctx context.Context, id int, opts ...RequestOption) (*model.BangumiSubject, error) {
	url := apiURL(APIPathGetSubject, id)
	subject := model.BangumiSubject{}

	req := c.client.R().
		SetContext(ctx).
		SetHeader("User-Agent", UserAgentHeader).
		SetHeader("Content-Type", ContentTypeJSON).
		SetResult(&subject).
		SetError(model.BangumiGenericErrorResponse{})

	applyRequestOptions(req, opts...)

	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, newAPIError(resp, url, ErrorGeneric)
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
		return nil, newAPIError(resp, url, ErrorGeneric)
	}

	return characters, nil
}

func (c *Client) RefreshAccessToken(ctx context.Context, token model.FirestoreBangumiToken) (*model.BangumiOAuthResponse, error) {
	tokenResp := model.BangumiOAuthResponse{}
	formData := map[string]string{
		"grant_type":    string(RefreshToken),
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
		return nil, newAPIError(resp, OAuthURL, ErrorOAuth)
	}

	return &tokenResp, nil
}

func (c *Client) GetHTML(ctx context.Context, path string) (*goquery.Document, error) {
	url := fmt.Sprintf("%s%s", HTMLBaseURL, path)

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

func batchFetch[T any](
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

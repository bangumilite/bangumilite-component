package bangumi

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/sstp105/bangumi-component/model"
	"net/http"
	"strings"
	"time"
)

var DefaultMaxRetries = 3
var DefaultInitialDelay = time.Duration(5) * time.Second
var DefaultMaxDelay = time.Duration(60) * time.Second

type Path string

const (
	Host string = "https://bangumi.tv"

	APIBaseURL    string = ""
	MonoPath      Path   = "mono"
	CharacterPath Path   = "character"
	PersonPath    Path   = "person"

	UserAgent                     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	ApplicationJsonContentType    = "application/json"
	ApplicationFormUrlencodedType = "application/x-www-form-urlencoded"
)

type ClientConfig struct {
	RetryCnt     int
	InitialDelay time.Duration
	MaxWaitTime  time.Duration
}

type Client struct {
	cfg    *ClientConfig
	client *resty.Client
	logger *logrus.Logger
}

func NewClient(logger *logrus.Logger) *Client {
	cfg := ClientConfig{
		RetryCnt:     DefaultMaxRetries,
		InitialDelay: DefaultInitialDelay,
		MaxWaitTime:  DefaultMaxDelay,
	}
	return NewWithConfig(logger, &cfg)
}

func NewWithConfig(logger *logrus.Logger, cfg *ClientConfig) *Client {
	if cfg == nil {
		return NewClient(logger)
	}

	if cfg.RetryCnt == 0 {
		cfg.RetryCnt = DefaultMaxRetries
	}

	if cfg.InitialDelay == 0 {
		cfg.InitialDelay = DefaultInitialDelay
	}

	if cfg.MaxWaitTime == 0 {
		cfg.MaxWaitTime = DefaultMaxDelay
	}

	client := resty.New()
	client.
		SetRetryCount(cfg.RetryCnt).
		SetRetryWaitTime(cfg.InitialDelay).
		SetRetryMaxWaitTime(cfg.MaxWaitTime).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() >= http.StatusInternalServerError
		})

	return &Client{
		cfg:    cfg,
		client: client,
		logger: logger,
	}
}

func (c *Client) GetSubject(ctx context.Context, id string, accessToken string) (*model.Subject, error) {
	url := fmt.Sprintf("https://api.bgm.tv/v0/subjects/%s", id)

	subject := model.Subject{}
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("User-Agent", UserAgent).
		SetHeader("Content-Type", ApplicationJsonContentType).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetResult(&subject).
		SetError(model.BangumiGenericErrorResponse{}).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		errResp := resp.Error().(*model.BangumiGenericErrorResponse)
		return nil, fmt.Errorf("failed to get subject %s, status code:%d, error:%s,%s",
			id,
			resp.StatusCode(),
			errResp.Title,
			errResp.Description,
		)
	}

	return &subject, nil
}

func (c *Client) RefreshAccessToken(ctx context.Context, token model.BangumiToken) (*model.BangumiTokenResponse, error) {
	formData := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     token.ClientID,
		"client_secret": token.ClientSecret,
		"redirect_uri":  token.RedirectURI,
		"refresh_token": token.RefreshToken,
	}

	tokenResp := model.BangumiTokenResponse{}

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", ApplicationFormUrlencodedType).
		SetFormData(formData).
		SetResult(&tokenResp).
		SetError(model.BangumiOAuthErrorResponse{}).
		Post("https://bgm.tv/oauth/access_token")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		errResp := resp.Error().(*model.BangumiOAuthErrorResponse)
		return nil, fmt.Errorf("failed to request new access token, status code:%d, message:%s",
			resp.StatusCode(),
			errResp.ErrorDescription,
		)
	}

	return &tokenResp, nil
}

func (c *Client) GetHTML(ctx context.Context, path Path) (*goquery.Document, error) {
	url := fmt.Sprintf("%s/%s", Host, path)

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

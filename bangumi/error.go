package bangumi

import (
	"fmt"
	"github.com/bangumilite/bangumilite-component/model"
	"github.com/go-resty/resty/v2"
)

func newAPIError(resp *resty.Response, path string, errorType APIError) error {
	switch errorType {
	case ErrorOAuth:
		errResp := resp.Error().(*model.BangumiOAuthErrorResponse)
		return apiError(path, resp.StatusCode(), errResp.Error, errResp.ErrorDescription)
	case ErrorGeneric:
		errResp := resp.Error().(*model.BangumiGenericErrorResponse)
		return apiError(path, resp.StatusCode(), errResp.Title, errResp.Description)
	default:
		return fmt.Errorf("unexpected error type: %s", errorType)
	}
}

func apiError(path string, statusCode int, error string, message string) error {
	return fmt.Errorf("failed to call %s, status code: %d, error: %s, message: %s",
		path,
		statusCode,
		error,
		message,
	)
}

package bangumi

import (
	"context"
	"fmt"
	"github.com/jarcoal/httpmock"
	_ "github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sstp105/bangumi-component/model"
	"net/http"
)

var _ = Describe("Bangumi API Unit Tests", func() {
	var (
		logger *logrus.Logger
		client *Client
	)

	BeforeEach(func() {
		logger = logrus.New()
		client = NewClient(logger)

		httpmock.ActivateNonDefault(client.client.GetClient())
	})

	Describe("GetSubject", func() {
		mockAccessToken := "<MOCK_ACCESS_TOKEN>"
		mockSubjectID := 1

		It("returns error if request returns error", func() {
			httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bgm.tv/v0/subjects/%d", mockSubjectID),
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(404, `
						{
                        	"title": "Bad Request",
						    "description": "Subject does not exist"
						}
					`,
					)
					resp.Header.Add("Content-Type", "application/json")
					return resp, nil
				},
			)

			resp, err := client.GetSubject(context.Background(), mockSubjectID, mockAccessToken)

			Expect(err).ToNot(BeNil())
			Expect(resp).To(BeNil())
		})

		It("returns error if cannot parse the response", func() {
			httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bgm.tv/v0/subjects/%d", mockSubjectID),
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, `
						{
						    "name": 123
						}
					`,
					)
					resp.Header.Add("Content-Type", "application/json")
					return resp, nil
				},
			)

			resp, err := client.GetSubject(context.Background(), mockSubjectID, mockAccessToken)

			Expect(err).ToNot(BeNil())
			Expect(resp).To(BeNil())
		})

		It("returns the subject if request succeed", func() {
			httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bgm.tv/v0/subjects/%d", mockSubjectID),
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, `
						{
					    	"id": 1,
					    	"name": "string",
					    	"name_cn": "string",
					    	"summary": "string"
						}
					`,
					)
					resp.Header.Add("Content-Type", "application/json")
					return resp, nil
				},
			)

			resp, err := client.GetSubject(context.Background(), mockSubjectID, mockAccessToken)

			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(resp.ID).To(Equal(1))
		})
	})

	Describe("GetSubjects", func() {
		mockAccessToken := "<MOCK_ACCESS_TOKEN>"
		mockSubjectIDs := []int{1, 2}

		It("returns the subject if request succeed", func() {
			httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bgm.tv/v0/subjects/%d", mockSubjectIDs[0]),
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, `
						{
					    	"id": 1,
					    	"name": "string",
					    	"name_cn": "string",
					    	"summary": "string"
						}
					`,
					)
					resp.Header.Add("Content-Type", "application/json")
					return resp, nil
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bgm.tv/v0/subjects/%d", mockSubjectIDs[1]),
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, `
						{
					    	"id": 2,
					    	"name": "string",
					    	"name_cn": "string",
					    	"summary": "string"
						}
					`,
					)
					resp.Header.Add("Content-Type", "application/json")
					return resp, nil
				},
			)

			resp, err := client.GetSubjects(context.Background(), mockSubjectIDs, mockAccessToken)

			Expect(err).To(BeNil())
			Expect(len(resp)).To(Equal(len(mockSubjectIDs)))
		})
	})

	Describe("RefreshAccessToken", func() {
		mockToken := model.BangumiToken{
			AccessToken:  "<ACCESS_TOKEN>",
			RefreshToken: "<REFRESH_TOKEN>",
			ClientID:     "<CLIENT_ID>",
			ClientSecret: "<CLIENT_SECRET>",
			RedirectURI:  "<REDIRECT_URI>",
		}

		It("returns error if request returns error", func() {
			httpmock.RegisterResponder("POST", "https://bgm.tv/oauth/access_token",
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(400, `
						{
							"error": "invalid_grant",
							"error_description": "Invalid refresh token"
						}
					`,
					)
					resp.Header.Add("Content-Type", "application/json")
					return resp, nil
				},
			)

			resp, err := client.RefreshAccessToken(context.Background(), mockToken)

			Expect(err).ToNot(BeNil())
			Expect(resp).To(BeNil())
		})

		It("returns new token if request succeeds", func() {
			httpmock.RegisterResponder("POST", "https://bgm.tv/oauth/access_token",
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, `
						{
							"access_token": "<NEW_ACCESS_TOKEN>",
							"refresh_token": "<NEW_REFRESH_TOKEN>",
							"expires_in": 100,
							"token_type": "<GRANT_TYPE>",
							"user_id": 1
						}
					`,
					)
					resp.Header.Add("Content-Type", "application/json")
					return resp, nil
				},
			)

			resp, err := client.RefreshAccessToken(context.Background(), mockToken)

			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(resp.AccessToken).To(Equal("<NEW_ACCESS_TOKEN>"))
			Expect(resp.RefreshToken).To(Equal("<NEW_REFRESH_TOKEN>"))
		})

		It("returns error if unable to parse the response", func() {
			httpmock.RegisterResponder("POST", "https://bgm.tv/oauth/access_token",
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, `
						{
							"access_token": "<NEW_ACCESS_TOKEN>",
							"refresh_token": "<NEW_REFRESH_TOKEN>",
							"expires_in": 100,
							"token_type": "<GRANT_TYPE>",
							"user_id": "<INVALID_USER_ID>"
						}
					`,
					)
					resp.Header.Add("Content-Type", "application/json")
					return resp, nil
				},
			)

			resp, err := client.RefreshAccessToken(context.Background(), mockToken)

			Expect(err).ToNot(BeNil())
			Expect(resp).To(BeNil())
		})
	})
})

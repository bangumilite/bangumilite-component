package model

import "strings"

type BangumiTags []BangumiTag

type BangumiSubject struct {
	ID         int               `json:"id" firestore:"id"`
	Type       int               `json:"type,omitempty" firestore:"type,omitempty"`
	Name       string            `json:"name" firestore:"name"`
	NameCn     string            `json:"name_cn" firestore:"name_cn"`
	Summary    string            `json:"summary" firestore:"summary"`
	Images     BangumiImages     `json:"images" firestore:"images"`
	Collection BangumiCollection `json:"collection" firestore:"collection"`
	Tags       BangumiTags       `json:"tags" firestore:"tags"`
	Rating     BangumiRating     `json:"rating" firestore:"rating"`
}

type BangumiImages struct {
	Small  string `json:"small" firestore:"small"`
	Medium string `json:"medium" firestore:"medium"`
	Large  string `json:"large" firestore:"large"`
}

type BangumiCollection struct {
	Wish    int `json:"wish" firestore:"wish"`
	Collect int `json:"collect" firestore:"collect"`
	Doing   int `json:"doing" firestore:"doing"`
	OnHold  int `json:"on_hold" firestore:"on_hold"`
	Dropped int `json:"dropped" firestore:"dropped"`
}

func (c BangumiCollection) Total() int {
	return c.Wish + c.Collect + c.Doing + c.OnHold + c.Dropped
}

type BangumiTag struct {
	Name string `json:"name" firestore:"name"`
}

func (t BangumiTags) ToString() string {
	var tags strings.Builder

	for i := 0; i < len(t) && i < 5; i++ {
		if i > 0 {
			tags.WriteString(" / ")
		}
		tags.WriteString(t[i].Name)
	}

	return tags.String()
}

type BangumiRating struct {
	Rank  int     `json:"rank" firestore:"rank"`
	Score float64 `json:"score" firestore:"score"`
}

type BangumiRelatedCharacter struct {
	Actors []BangumiPerson `json:"actors" firestore:"actors"`
}

type BangumiPerson struct {
	ID   int    `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
}

type BangumiOAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	UserID       int    `json:"user_id"`
}

type BangumiOAuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type BangumiGenericErrorResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type SubjectTypeID string

const (
	BookID  SubjectTypeID = "1"
	AnimeID SubjectTypeID = "2"
	MusicID SubjectTypeID = "3"
	GameID  SubjectTypeID = "4"
	RealID  SubjectTypeID = "6"
)

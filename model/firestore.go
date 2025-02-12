package model

import (
	"errors"
)

type FirestoreBangumiToken struct {
	AccessToken  string `firestore:"access_token"`
	RefreshToken string `firestore:"refresh_token"`
	ClientID     string `firestore:"client_id"`
	ClientSecret string `firestore:"client_secret"`
	RedirectURI  string `firestore:"redirect_uri"`
}

type FirestoreSubject struct {
	ID         int     `firestore:"id" json:"id"`
	Name       string  `firestore:"name" json:"name"`
	NameCn     string  `firestore:"name_cn" json:"name_cn"`
	Image      string  `firestore:"image,omitempty" json:"image"`
	Info       string  `firestore:"info" json:"info"`
	Score      float64 `firestore:"score" json:"score"`
	Rank       int     `firestore:"rank,omitempty" json:"rank"`
	Collection int     `firestore:"collection" json:"collection"`
}

type FirestoreSeasonSubject struct {
	ID       int                `firestore:"id" json:"id"`
	Name     string             `firestore:"name" json:"name"`
	NameCn   string             `firestore:"name_cn" json:"name_cn"`
	Image    string             `firestore:"image" json:"image"`
	Info     string             `firestore:"info" json:"info"`
	Summary  string             `firestore:"summary" json:"summary"`
	Actors   []BangumiPerson    `firestore:"actors" json:"actors"`
	Staff    []string           `firestore:"staff" json:"staff"`
	Trailers []FirestoreTrailer `firestore:"trailers" json:"trailers"`
}

type FirestoreTrailer struct {
	URL       string `json:"url" firestore:"url"`
	Thumbnail string `json:"thumbnail" firestore:"thumbnail"`
}

type FirestoreSeasonIndexDocument struct {
	Data []FirestoreSeasonIndexItem `firestore:"data" json:"data"`
}

type FirestoreSeasonIndexItem struct {
	ID    string `firestore:"id" json:"id"`
	Image string `firestore:"image" json:"image"`
}

type FirestoreDiscoverySubject struct {
	Title string             `firestore:"title" json:"title"`
	Data  []FirestoreSubject `firestore:"data" json:"data"`
}

type FirestoreMonoDocument struct {
	Trending  []FirestoreMono `json:"trending" firestore:"trending"`
	Popular   []FirestoreMono `json:"popular" firestore:"popular"`
	Birthday  []FirestoreMono `json:"birthday" firestore:"birthday"`
	Inventory []FirestoreMono `json:"inventory" firestore:"inventory"`
}

type FirestoreMono struct {
	ID                int                         `json:"id" firestore:"id"`
	Name              string                      `json:"name" firestore:"name"`
	Image             *string                     `json:"image" firestore:"image,omitempty"`
	RelatedSubjects   *[]FirestoreMonoRelatedWork `json:"related_subjects,omitempty" firestore:"related_subjects,omitempty"`
	RelatedCharacters *[]FirestoreMonoRelatedWork `json:"related_characters,omitempty" firestore:"related_characters,omitempty"`
}

type FirestoreMonoRelatedWork struct {
	ID       int            `json:"id" firestore:"id"`
	Name     string         `json:"name" firestore:"name"`
	Type     *int           `json:"type" firestore:"type"`
	NameCN   *string        `json:"name_cn" firestore:"name_cn,omitempty"`
	Relation *string        `json:"relation" firestore:"relation,omitempty"`
	Image    *string        `json:"image" firestore:"image,omitempty"`
	Mono     *FirestoreMono `json:"mono" firestore:"mono,omitempty"`
}

func (m FirestoreMonoDocument) Validate() error {
	if len(m.Trending) == 0 {
		return errors.New("trending mono is empty")
	}

	if len(m.Popular) == 0 {
		return errors.New("popular mono is empty")
	}

	if len(m.Birthday) == 0 {
		return errors.New("birthday mono is empty")
	}

	if len(m.Inventory) == 0 {
		return errors.New("inventory mono is empty")
	}

	return nil
}

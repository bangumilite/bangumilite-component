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
	// Trending represents bangumi.tv current hot mono, calculated based on recent 30 days collections.
	Trending []FirestoreMono `json:"trending" firestore:"trending"`

	// Popular represents the most popular mono all time.
	Popular []FirestoreMono `json:"popular" firestore:"popular"`

	// Birthday represents the today's birthday mono (Updated every day).
	Birthday []FirestoreMono `json:"birthday" firestore:"birthday"`

	// Inventory represents the newly added mono in bangumi.
	Inventory []FirestoreMono `json:"inventory" firestore:"inventory"`
}

// FirestoreMono represents a character or person type.
type FirestoreMono struct {
	// ID is person or character id.
	ID int `json:"id" firestore:"id"`

	// Name is mono name in original language.
	Name string `json:"name" firestore:"name"`

	// Image is the bangumi mono grid image, which is equivalent to bangumi image.grid.
	Image *string `json:"image" firestore:"image,omitempty"`

	// RelatedSubjects holds an array mono related subjects.
	// For character, it includes the subjects that the character recently acted.
	// For person, it can be the subjects that the person recently participated.
	RelatedSubjects *[]FirestoreMonoRelatedWork `json:"related_subjects,omitempty" firestore:"related_subjects,omitempty"`

	// RelatedCharacters is an array of MonoRelatedWork.
	// For voice actor, it includes the characters that the person recently starred.
	RelatedCharacters *[]FirestoreMonoRelatedWork `json:"related_characters,omitempty" firestore:"related_characters,omitempty"`
}

// FirestoreMonoRelatedWork represents a related work (subject) for a Mono (person or character).
// For person, the related work could be 最近参与. The Mono field will always be nil.
// For person voice actor, the work can include 最近演出角色. The Mono field will represent the character.
// For character, it can only be 最近演出 and the Mono field will refer to the voice actor.
type FirestoreMonoRelatedWork struct {
	// ID is bangumi subject_id.
	ID int `json:"id" firestore:"id"`

	// Name is subject title in original language.
	Name string `json:"name" firestore:"name"`

	// Type is bangumi subject_type (1,2,3,4,6)
	Type *int `json:"type" firestore:"type"`

	// NameCN is subject title in chinese, can be nil.
	NameCN *string `json:"name_cn" firestore:"name_cn,omitempty"`

	// Relation is either the relation between the character and the subject, or the person to the subject.
	// E.g. for character, it can be the 主角,配角...
	// For person, it can be 艺术家,脚本...
	Relation *string `json:"relation" firestore:"relation,omitempty"`

	// Image is the subject image in high resolution, equivalent to bangumi image.large.
	Image *string `json:"image" firestore:"image,omitempty"`

	// Mono is the related mono. It could be character or voice actor.
	Mono *FirestoreMono `json:"mono" firestore:"mono,omitempty"`
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

package model

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

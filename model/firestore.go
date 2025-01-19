package model

type BangumiToken struct {
	AccessToken  string `firestore:"access_token"`
	RefreshToken string `firestore:"refresh_token"`
	ClientID     string `firestore:"client_id"`
	ClientSecret string `firestore:"client_secret"`
	RedirectURI  string `firestore:"redirect_uri"`
}

type TrendingSubject struct {
	ID         int     `firestore:"id" json:"id"`
	Name       string  `firestore:"name" json:"name"`
	NameCn     string  `firestore:"name_cn" json:"name_cn"`
	Image      string  `firestore:"image" json:"image"`
	Info       string  `firestore:"info" json:"info"`
	Score      float64 `firestore:"score" json:"score"`
	Rank       *int    `firestore:"rank" json:"rank"`
	Collection int     `firestore:"collection" json:"collection"`
}

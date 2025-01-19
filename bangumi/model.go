package bangumi

type Subject struct {
	ID         int        `json:"id" firestore:"id"`
	Name       string     `json:"name" firestore:"name"`
	NameCn     string     `json:"name_cn" firestore:"name_cn"`
	Summary    string     `json:"summary" firestore:"summary"`
	Images     Images     `json:"images" firestore:"images"`
	Collection Collection `json:"collection" firestore:"collection"`
	Tags       []Tag      `json:"tags" firestore:"tags"`
	Rating     Rating     `json:"rating" firestore:"rating"`
}

type Images struct {
	Small  string `json:"small" firestore:"small"`
	Medium string `json:"medium" firestore:"medium"`
	Large  string `json:"large" firestore:"large"`
}

type Collection struct {
	Wish    int `json:"wish" firestore:"wish"`
	Collect int `json:"collect" firestore:"collect"`
	Doing   int `json:"doing" firestore:"doing"`
	OnHold  int `json:"on_hold" firestore:"on_hold"`
	Dropped int `json:"dropped" firestore:"dropped"`
}

func (c Collection) Total() int {
	return c.Wish + c.Collect + c.Doing + c.OnHold + c.Dropped
}

type Tag struct {
	Name string `json:"name" firestore:"name"`
}
type Rating struct {
	Rank  int `json:"rank" firestore:"rank"`
	Score float64
}

type RelatedCharacter struct {
	Actors []Person `json:"actors" firestore:"actors"`
}

type Person struct {
	ID   int    `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
}

type OAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	UserID       int    `json:"user_id"`
}

type OAuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type GenericErrorResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

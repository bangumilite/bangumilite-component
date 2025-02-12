package fs

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/bangumilite/bangumilite-component/mailer"
	"github.com/bangumilite/bangumilite-component/model"
	"google.golang.org/api/option"
	"os"
)

const (
	FirebaseProjectID = "yang-bangumi"

	TokenCollectionKey           = "token"
	TokenCollectionBangumiDocKey = "bangumi"

	SeasonCollectionKey         = "season"
	SeasonCollectionIndexDocKey = "index"

	MonoCollectionKey = "mono"

	DiscoveryCollectionKey = "discovery"

	MailgunDocumentKey = "mailgun"

	BangumiAccessTokenKey  = "access_token"
	BangumiRefreshTokenKey = "refresh_token"

	FirebaseLastUpdatedTimestampKey = "lastUpdatedDate"
)

var ErrDocumentDoesNotExist = errors.New("document does not exist")

type Client struct {
	fs *firestore.Client
}

func New(ctx context.Context) (*Client, error) {
	var fs *firestore.Client
	var err error

	env := os.Getenv(model.RunningEnvironment)
	switch env {
	case string(model.Production):
		fs, err = firestore.NewClient(ctx, FirebaseProjectID)
	default:
		fs, err = firestore.NewClient(ctx, FirebaseProjectID, option.WithCredentialsFile("service_account.json"))
	}

	if err != nil {
		return nil, err
	}

	return &Client{
		fs: fs,
	}, nil
}

func (c *Client) Close() error {
	return c.fs.Close()
}

func (c *Client) GetBangumiToken(ctx context.Context) (*model.FirestoreBangumiToken, error) {
	docRef := c.fs.Collection(TokenCollectionKey).Doc(TokenCollectionBangumiDocKey)

	data, err := getDocument[model.FirestoreBangumiToken](ctx, docRef)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) GetMailgunConfig(ctx context.Context) (*mailer.MailgunConfig, error) {
	docRef := c.fs.Collection(TokenCollectionKey).Doc(MailgunDocumentKey)

	data, err := getDocument[mailer.MailgunConfig](ctx, docRef)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) GetSeasonIndex(ctx context.Context) (*model.FirestoreSeasonIndexDocument, error) {
	docRef := c.fs.Collection(SeasonCollectionKey).Doc(SeasonCollectionIndexDocKey)

	data, err := getDocument[model.FirestoreSeasonIndexDocument](ctx, docRef)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) UpdateMonoDocument(ctx context.Context, monoType model.MonoType, data model.FirestoreMonoDocument) error {
	docRef := c.fs.Collection(MonoCollectionKey).Doc(string(monoType))
	docData := map[string]interface{}{
		"data":                          data,
		FirebaseLastUpdatedTimestampKey: firestore.ServerTimestamp,
	}

	err := saveDocument(ctx, docRef, docData)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateSeasonIndex(ctx context.Context, items []model.FirestoreSeasonIndexItem) error {
	docRef := c.fs.Collection(SeasonCollectionKey).Doc(SeasonCollectionIndexDocKey)

	data := map[string]interface{}{
		"data":                          items,
		FirebaseLastUpdatedTimestampKey: firestore.ServerTimestamp,
	}

	err := saveDocument(ctx, docRef, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateBangumiToken(ctx context.Context, accessToken string, refreshToken string) error {
	docRef := c.fs.Collection(TokenCollectionKey).Doc(TokenCollectionBangumiDocKey)

	data := map[string]interface{}{
		BangumiAccessTokenKey:           accessToken,
		BangumiRefreshTokenKey:          refreshToken,
		FirebaseLastUpdatedTimestampKey: firestore.ServerTimestamp,
	}

	err := saveDocument(ctx, docRef, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateTrendingSubjects(ctx context.Context, subjectTypeID string, subjects []model.FirestoreSubject) error {
	docRef := c.fs.Collection("trending").Doc(subjectTypeID)

	data := map[string]interface{}{
		"data":                          subjects,
		FirebaseLastUpdatedTimestampKey: firestore.ServerTimestamp,
	}

	err := saveDocument(ctx, docRef, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateSeasonalSubjects(ctx context.Context, id string, subjects []model.FirestoreSeasonSubject) error {
	docRef := c.fs.Collection("season").Doc(id)

	data := map[string]interface{}{
		"data":                          subjects,
		"total":                         len(subjects),
		FirebaseLastUpdatedTimestampKey: firestore.ServerTimestamp,
	}

	err := saveDocument(ctx, docRef, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateDiscoverySubjects(ctx context.Context, id model.SubjectTypeID, data []model.FirestoreDiscoverySubject) error {
	docRef := c.fs.Collection(DiscoveryCollectionKey).Doc(string(id))

	docData := map[string]interface{}{
		"data":                          data,
		FirebaseLastUpdatedTimestampKey: firestore.ServerTimestamp,
	}

	err := saveDocument(ctx, docRef, docData)
	if err != nil {
		return err
	}

	return nil
}

func getDocument[T any](ctx context.Context, docRef *firestore.DocumentRef) (*T, error) {
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	if !docSnap.Exists() {
		return nil, ErrDocumentDoesNotExist
	}

	var result T
	err = docSnap.DataTo(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func saveDocument(ctx context.Context, docRef *firestore.DocumentRef, data map[string]interface{}) error {
	_, err := docRef.Set(ctx, data, firestore.MergeAll)

	if err != nil {
		return err
	}

	return nil
}

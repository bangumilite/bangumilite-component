package fs

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sstp105/bangumi-component/mailer"
	"github.com/sstp105/bangumi-component/model"
	"google.golang.org/api/option"
	"os"
)

const (
	FirebaseProjectID = "yang-bangumi"

	TokenCollectionKey           = "token"
	TokenCollectionBangumiDocKey = "bangumi"

	SeasonCollectionKey         = "season"
	SeasonCollectionIndexDocKey = "index"

	MailgunDocumentKey = "mailgun"

	BangumiAccessTokenKey  = "access_token"
	BangumiRefreshTokenKey = "refresh_token"

	FirebaseLastUpdatedTimestampKey = "lastUpdatedDate"
)

// Client is a wrapper around the Firestore client and providing additional logging
type Client struct {
	fs     *firestore.Client
	logger *logrus.Logger
}

var ErrDocumentDoesNotExist = errors.New("document does not exist")

// New creates a new Firestore client.
func New(ctx context.Context, logger *logrus.Logger, option ...option.ClientOption) (*Client, error) {
	var fs *firestore.Client
	var err error

	env := os.Getenv(model.RunningEnvironment)
	switch env {
	case string(model.Production):
		fs, err = firestore.NewClient(ctx, FirebaseProjectID)
	default:
		fs, err = firestore.NewClient(ctx, FirebaseProjectID, option...)
	}

	if err != nil {
		return nil, err
	}

	return &Client{
		fs:     fs,
		logger: logger,
	}, nil
}

// Close gives caller access to close the firestore client.
func (c *Client) Close() error {
	return c.fs.Close()
}

// GetBangumiToken retrieves bangumi tokens that can be used to call bangumi API.
func (c *Client) GetBangumiToken(ctx context.Context) (*model.FirestoreBangumiToken, error) {
	docRef := c.fs.Collection(TokenCollectionKey).Doc(TokenCollectionBangumiDocKey)

	data, err := getDocument[model.FirestoreBangumiToken](ctx, docRef)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetMailgunConfig retrieves the Mailgun related configuration from firestore.
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

func (c *Client) UpdateSeasonIndex(ctx context.Context, items []model.FirestoreSeasonIndexItem) error {
	docRef := c.fs.Collection(SeasonCollectionKey).Doc(SeasonCollectionIndexDocKey)

	data := map[string]interface{}{
		"data":            items,
		"lastUpdatedDate": firestore.ServerTimestamp,
	}

	err := saveDocument(ctx, docRef, data)
	if err != nil {
		return err
	}

	return nil
}

// UpdateBangumiToken overrides access_token and refresh_token with new valid tokens.
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

func (c *Client) UpdateTrendingSubjects(ctx context.Context, subjectTypeID string, subjects []model.FirestoreTrendingSubject) error {
	docRef := c.fs.Collection("trending").Doc(subjectTypeID)

	data := map[string]interface{}{
		"data":            subjects,
		"lastUpdatedDate": firestore.ServerTimestamp,
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
		"data":  subjects,
		"total": len(subjects),
	}

	err := saveDocument(ctx, docRef, data)
	if err != nil {
		return err
	}

	return nil
}

// getDocument retrieves a Firestore document and unmarshal its data into a specified type.
//
// This function is generic and can be used to fetch and unmarshal documents into any struct type.
//
// Type Parameters:
//   - T: The type into which the Firestore document data will be unmarshalled.
//
// Parameters:
//   - ctx: The context for the Firestore operation.
//   - docRef: A reference to the Firestore document to retrieve.
//
// Returns:
//   - *T: A pointer to the unmarshalled result of type T, or nil if an error occurs.
//   - error: An error if the operation fails or if the document does not exist.
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

// saveDocument saves data to a Firestore document.
//
// This function updates the specified Firestore document with the provided data. If the document
// does not exist, it will be created. The operation merges the provided data with any existing
// fields in the document.
//
// Parameters:
//   - ctx: The context for the Firestore operation.
//   - docRef: A reference to the Firestore document to save or update.
//   - data: A map of key-value pairs representing the data to save in the document.
//
// Returns:
//   - error: An error if the operation fails, or nil if the operation is successful.
//
// Firestore Merge Behavior:
//   - The `firestore.MergeAll` option ensures that only the fields in the `data` map
//     are updated or added to the document. Existing fields not included in the `data` map
//     remain unchanged.
func saveDocument(ctx context.Context, docRef *firestore.DocumentRef, data map[string]interface{}) error {
	_, err := docRef.Set(ctx, data, firestore.MergeAll)

	if err != nil {
		return err
	}

	return nil
}

package functions

// Function triggered when thumbnail image uploaded to cloud storage bucket
// Update user document in firestore with location of thumbnail
//
// Following env variables required:
//   GCLOUD_PROJECT_ID - used to init firestore client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	errors_pkg "github.com/pkg/errors"

	"encoding/json"

	"cloud.google.com/go/firestore"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// UpdateFirebaseSentiment update firebase user doc with sentiment
func UpdateFirebaseSentiment(ctx context.Context, msg PubSubMessage) error {

	var sentiment Sentiment
	err := json.Unmarshal(msg.Data, &sentiment)
	if err != nil {
		return errors_pkg.Wrap(err, "could not convert bytes to sentiment")
	}

	storageImage, err := ToStorageImage(sentiment.GCSEvent)
	if err != nil {
		return err
	}

	log.Printf("processing pubsub event for storage event %v", sentiment.GCSEvent)

	projectID := os.Getenv("GCLOUD_PROJECT_ID")
	if projectID == "" {
		return errors.New("set (Firebase) project ID via GCLOUD_PROJECT_ID env variable")
	}
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("cannot create firebase client: %v", err)
	}
	defer client.Close()

	photoDoc := client.Collection("users").Doc(storageImage.UserID).Collection("photos").Doc(storageImage.PhotoID())
	_, err = photoDoc.Update(ctx, []firestore.Update{
		{
			Path:  "sentiment",
			Value: sentiment.Likelihoods,
		},
	})
	if err != nil {
		return fmt.Errorf("error updating firebase: %s", err)
	}

	log.Printf("updated firstore document users/%s/photos/%s value 'sentiment' with: %s", storageImage.UserID, storageImage.PhotoID(), sentiment.Likelihoods)

	return nil
}

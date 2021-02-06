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

	"cloud.google.com/go/firestore"
)

// UpdateFirebaseThumb update firebase user doc with location of thumbnail
func UpdateFirebaseThumb(ctx context.Context, event GCSEvent) error {
	log.Printf("processing gs://%s/%s", event.Bucket, event.Name)

	storageImage, err := ToStorageImage(event)
	if err != nil {
		return err
	}

	// if its not a thumbnail, ignore
	if !storageImage.IsThumbNail() {
		log.Printf("skipping processing non thumbnail image (%s)", storageImage.ImageName)
		return nil
	}

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
	thubURI := event.Name
	_, err = photoDoc.Update(ctx, []firestore.Update{
		{
			Path:  "thumbUri",
			Value: thubURI,
		},
	})
	if err != nil {
		return fmt.Errorf("error updating firebase: %s", err)
	}

	log.Printf("updated firstore document users/%s/photos/%s value 'thumbUri' with: %s", storageImage.UserID, storageImage.PhotoID(), thubURI)

	return nil
}

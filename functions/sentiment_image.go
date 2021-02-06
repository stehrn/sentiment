package functions

// Required Env variables (used to publish Sentiment to pub/sub topic)

// Following env variables required:
//   GCLOUD_PROJECT_ID - used to init pubsub client
//   RESULT_TOPIC - topic photo sentiment published to

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"github.com/pkg/errors"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// Likelihoods likelihood of different sentiments
type Likelihoods struct {
	JoyLikelihood      visionpb.Likelihood `json:"JoyLikelihood"`
	SorrowLikelihood   visionpb.Likelihood `json:"SorrowLikelihood"`
	AngerLikelihood    visionpb.Likelihood `json:"AngerLikelihood"`
	SurpriseLikelihood visionpb.Likelihood `json:"SurpriseLikelihood"`
}

// Sentiment sentiment of an image
type Sentiment struct {
	// GCSEvent cloud storage info
	GCSEvent GCSEvent `json:"GCSEvent"`
	// Likelihood likelihood of image sentiments
	Likelihoods Likelihoods `json:"Likelihoods"`
}

// Marshal conver Sentiment to []byte
func (sentiment Sentiment) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(sentiment)
	if err != nil {
		return nil, errors.Wrap(err, "Could not convert sentiment to bytes")
	}
	return bytes, nil
}

// ProcessImageSentiment download image from cloud storage and use cloud vision vision to derive image sentiment,
// publish sentiment to pub/sub topic
func ProcessImageSentiment(ctx context.Context, event GCSEvent) error {
	log.Printf("Processing gs://%s/%s", event.Bucket, event.Name)

	if event.Bucket == "" {
		return fmt.Errorf("empty event.Bucket")
	}
	if event.Name == "" {
		return fmt.Errorf("empty event.Name")
	}

	storageImage, err := ToStorageImage(event)
	if err != nil {
		return err
	}

	// dont process images we're already proceccessed
	if storageImage.IsThumbNail() {
		log.Printf("Not processing thumbnail (%s)", storageImage.ImageName)
		return nil
	}

	// read photo from cloud storage
	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer storageClient.Close()

	inputBlob := storageClient.Bucket(event.Bucket).Object(event.Name)
	photo, err := inputBlob.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("NewReader: %v", err)
	}

	// work out likelihood of different sentiments
	likelihoods, err := sentiment(photo)
	if err != nil {
		return err
	}

	// publish sentiment
	sentiment := Sentiment{
		GCSEvent:    event,
		Likelihoods: *likelihoods,
	}
	if err := publish(ctx, sentiment); err != nil {
		return fmt.Errorf("Error publishing sentiment for event: %v, error: %v", event, err)
	}

	log.Printf("photo %s processed", event.Name)
	return nil
}

// Sentiment apply AI to work out sentiment of face detected in image
func sentiment(file *storage.Reader) (*Likelihoods, error) {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("vision.NewImageAnnotatorClient: %v", err)
	}
	defer client.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return nil, fmt.Errorf("vision.NewImageFromReader: %v", err)
	}
	annotations, err := client.DetectFaces(ctx, image, nil, 10)
	if err != nil {
		return nil, fmt.Errorf("client.DetectFaces: %v", err)
	}

	var likelihood Likelihoods
	if len(annotations) != 0 {
		if len(annotations) > 1 {
			log.Printf("Found %d faces, just picking first one", len(annotations))
		}
		annotation := annotations[0]
		likelihood = Likelihoods{
			annotation.JoyLikelihood,
			annotation.SorrowLikelihood,
			annotation.AngerLikelihood,
			annotation.SurpriseLikelihood,
		}
	}
	return &likelihood, nil
}

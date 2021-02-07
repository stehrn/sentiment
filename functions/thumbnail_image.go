package functions

// Function triggered when image uploaded to cloud storage bucket
// Create thumbnail of image and store in same cloud storage
// location appended with '_thumb'

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"cloud.google.com/go/storage"
)

// GenerateThumbnailImage generate thumbnail from image stored
// at gs://[bucket]/[userId]/[imageName]
// and store result in gs://[bucket]/[userId]/thumb_[imageName]
func GenerateThumbnailImage(ctx context.Context, event GCSEvent) error {
	log.Printf("processing gs://%s/%s", event.Bucket, event.Name)

	storageImage, err := ToStorageImage(event)
	if err != nil {
		return err
	}

	// dont process images we're already proceccessed
	if storageImage.IsThumbNail() {
		log.Printf("skipping existing thumbnail (%s)", storageImage.ImageName)
		return nil
	}

	// create reader/writer streams frm/to cloud storage
	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer storageClient.Close()

	inputBlob := storageClient.Bucket(event.Bucket).Object(event.Name)
	r, err := inputBlob.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("inputBlob.NewReader: %v", err)
	}

	outputName := storageImage.ToThumbNail()
	outputBlob := storageClient.Bucket(event.Bucket).Object(outputName)

	w := outputBlob.NewWriter(ctx)
	defer w.Close()

	// convert image
	cmd := exec.Command("convert", "-", "-thumbnail", "200x200>", "-")
	cmd.Stdin = r
	cmd.Stdout = w

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd.Run: %v", err)
	}

	log.Printf("thumbnail uploaded to gs://%s/%s", outputBlob.BucketName(), outputBlob.ObjectName())
	return nil
}

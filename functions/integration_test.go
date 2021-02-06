package functions

// Test needs `convert` on filepath otherwise get this error:
//      ThumbnailImage cmd.Run: exec: "convert": executable file not found in $PATH

// used by integration test:
// export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/integration_test_key.json
// export CLOUD_STORAGE_BUCKET_NAME=int-bucket

// used by publish.go:
// export RESULT_TOPIC=sentiment-topic
// export GOOGLE_CLOUD_PROJECT=sentiment-302320

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/storage"
)

const object = "123456/photo.jpg"

var (
	bucket        = os.Getenv("CLOUD_STORAGE_BUCKET_NAME")
	storageClient *storage.Client
)

func TestMain(m *testing.M) {
	err := setUp()
	var code int
	if err != nil {
		code = -1
		log.Printf("Error during setUp: %v", err)
	} else {
		code = m.Run()
	}
	err = tearDown()
	if err != nil {
		code = -2
		log.Printf("Error during tearDown: %v", err)
	}
	os.Exit(code)
}

// Ignore for now as wont work without local install of 'convert' app
func IgnoreTestGenerateThumbnailImage(t *testing.T) {
	event := GCSEvent{
		Bucket: bucket,
		Name:   object,
	}

	ctx := context.Background()
	err := GenerateThumbnailImage(ctx, event)
	if err != nil {
		t.Error("Error processing ThumbnailImage", err)
	}
}

func TestProcessImageSentiment(t *testing.T) {
	event := GCSEvent{
		Bucket: bucket,
		Name:   object,
	}

	ctx := context.Background()
	err := ProcessImageSentiment(ctx, event)
	if err != nil {
		t.Error("error processing image sentiment", err)
	}
}

//upload photo to cloud storage
func setUp() error {
	var err error

	storageClient, err = storage.NewClient(context.Background())
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}

	content, err := ioutil.ReadFile("../img/photo.jpg")
	if err != nil {
		return err
	}

	ctx := context.Background()
	wc := storageClient.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, bytes.NewReader(content)); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func tearDown() error {
	if storageClient == nil {
		return nil
	}

	defer storageClient.Close()
	ctx := context.Background()
	err := storageClient.Bucket(bucket).Object(object).Delete(ctx)
	if err != nil {
		return fmt.Errorf("storageClient.Delete: %v", err)
	}
	return nil
}

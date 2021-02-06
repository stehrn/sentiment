package functions

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

// Marshaler marchalljson
type Marshaler interface {
	Marshal() ([]byte, error)
}

func publish(ctx context.Context, data Marshaler) error {
	var err error

	topicName := os.Getenv("RESULT_TOPIC")
	if topicName == "" {
		return errors.New("set result pubsub topic via RESULT_TOPIC env variable")
	}

	projectID := os.Getenv("GCLOUD_PROJECT_ID")
	if projectID == "" {
		return errors.New("set cloud project via GCLOUD_PROJECT_ID env variable")
	}

	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer pubsubClient.Close()

	// check topic exists
	topic := pubsubClient.Topic(topicName)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("topic.Exists: %v", err)
	}
	if !ok {
		return fmt.Errorf("topic %s does not exist", topicName)
	}

	bytes, err := data.Marshal()
	if err != nil {
		return err
	}

	res := topic.Publish(ctx, &pubsub.Message{Data: bytes})
	_, err = res.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed tp publish to topic %s, error: %v", topicName, err)
	}

	log.Printf("data published to topic %s", topicName)

	return nil
}

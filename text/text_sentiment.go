package main

import (
	"context"
	"log"

	lang "cloud.google.com/go/language/apiv1"
	langpb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

var (
	appContext context.Context
	langClient *lang.Client
)

func init() {
	appContext = context.Background()
	client, err := lang.NewClient(appContext)
	if err != nil {
		log.Panicf("Failed to create client: %v", err)
	}
	langClient = client
}

func scoreSentiment(s string) (*langpb.Sentiment, error) {
	result, err := langClient.AnalyzeSentiment(appContext, &langpb.AnalyzeSentimentRequest{
		Document: &langpb.Document{
			Source: &langpb.Document_Content{
				Content: s,
			},
			Type: langpb.Document_PLAIN_TEXT,
		},
		EncodingType: langpb.EncodingType_UTF8,
	})
	if err != nil {
		return nil, err
	}
	return result.DocumentSentiment, nil
}


## Overview
Sentiment analysis using https://cloud.google.com/natural-language/

To run from command line:
```
gcloud ml language analyze-entity-sentiment --content="I love R&B music. Marvin Gaye is the best. 'What's Going On' is one of my favorite songs. It was so sad when Marvin Gaye died."
```

Code is based on https://cloud.google.com/functions/docs/writing/http#writing_http_content-go, and consists of:
* [text_http.go](text_http.go) - Simple http server
* [text_sentiment.go](text_sentiment.go) - sentiment analysis using Googles golang natural language client API

## Deploy Cloud Functions
Open terminal at location of go scripts (and this README) and run following commands...

```
gcloud config set project <PROJECT_ID>
gcloud functions deploy sentiment_http --runtime=go113 --entry-point Analyse --trigger-http --allow-unauthenticated
```
(see https://cloud.google.com/sdk/gcloud/reference/functions/deploy)


Test the function
```
export HTTP_URL=$(gcloud functions describe sentiment_http --format="value(httpsTrigger.url)")
curl -X POST $HTTP_URL -H "Content-Type:application/json" --data '{"message":"I love to code"}'
```
You should get back a message showing sentiment:
```
magnitude:0.9  score:0.9
```
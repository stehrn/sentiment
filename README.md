# Overview
Firebase web app with several cloud functions running on GCP

* GitHub: https://github.com/stehrn/sentiment
* Firebase app: https://console.firebase.google.com/project/sentiment-302320/
* stackblitz: https://stackblitz.com/edit/sentiments

# Functions
User uploads a photo of their face to cloud storage at gs://bucket/userID/photoID.png, storage event triggers functions to generate thumbnail and determine image sentiment.

## Thumbnail
* Thumbnail function generates thumbnail and stores in cloud storage (gs://bucket/userID/thumb_photoID.png)
* Storage event from thumbnail triggers function to update firebase photo document at /users/userID/photos/photoID with location of thumbnail 

## Sentiment
* Sentiment function uses Cloud Vision API to determine sentiment, publishing result to pubsub topic
* Function subscribed to topic updates firebase photo document at /users/userID/photos/photoID with image sentiment

# Web App
Web app developed using [stackblitz](https://stackblitz.com/edit/sentiments)

# Deployment
Swtich to correct project:
```
export GCLOUD_PROJECT_ID=sentiment-302320
gcloud config set project $GCLOUD_PROJECT_ID
```

## Thumbnail Functions
```
export IMAGE_BUCKET_NAME=sentiment-302320.appspot.com
export REGION=europe-west2

gcloud functions deploy thumb-image \
--runtime go113 \
--region ${REGION} \
--trigger-bucket ${IMAGE_BUCKET_NAME} \
--entry-point GenerateThumbnailImage 

gcloud functions deploy thumb-firebase \
--runtime go113 \
--region ${REGION} \
--trigger-bucket ${IMAGE_BUCKET_NAME} \
--entry-point UpdateFirebaseThumb \
--set-env-vars "GCLOUD_PROJECT_ID=${GCLOUD_PROJECT_ID}"
```

## Sentiment Analysis Functions
```
export RESULT_TOPIC=sentiment-topic
gcloud pubsub topics create ${RESULT_TOPIC}

gcloud functions deploy sentiment-image \
--runtime go113 \
--region ${REGION} \
--trigger-bucket ${IMAGE_BUCKET_NAME} \
--entry-point ProcessImageSentiment \
--set-env-vars "^:^GCLOUD_PROJECT_ID=${GCLOUD_PROJECT_ID}:RESULT_TOPIC=${RESULT_TOPIC}"

gcloud functions deploy sentiment-firebase \
--runtime go113 \
--region ${REGION} \
--trigger-topic ${RESULT_TOPIC} \
--entry-point UpdateFirebaseSentiment \
--set-env-vars "GCLOUD_PROJECT_ID=${GCLOUD_PROJECT_ID}"
```

# Integration tests

## Set up service account permissions and generate key for use with integration test
One off tasks

Create service account:
```
gcloud iam service-accounts create sentiment-client --description="Sentiment client account" --display-name="Sentiment client account"

export SERVICE_ACCOUNT=sentiment-client@${PROJECT_ID}.iam.gserviceaccount.com
gcloud iam service-accounts describe ${SERVICE_ACCOUNT}
```
Bind role to SA:
```
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member=serviceAccount:${SERVICE_ACCOUNT} --role "roles/owner"
```
(yes, role very broad, can look to limit scope)

Create key:
```
gcloud iam service-accounts keys create ${HOME}/sentiment_key.json --iam-account ${SERVICE_ACCOUNT}
```
Export into env:
```
export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/sentiment_key.json
```

## Run tests

```
// used by integration test:
export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/integration_test_key.json
export CLOUD_STORAGE_BUCKET_NAME=int-bucket

// used by publish.go:
export RESULT_TOPIC=sentiment-topic 
export GOOGLE_CLOUD_PROJECT=sentiment-302320 

go test -v
```

## Ad-hoc commands 

```
// create topic, subscription, and pull message
gcloud pubsub topics list
gcloud pubsub topics create sentiment-topic
gcloud pubsub subscriptions create tmp --topic projects/sentiment-302320/topics/sentiment-topic
gcloud pubsub subscriptions pull tmp
```

## Reading 

* https://cloud.google.com/functions/docs/tutorials/ocr
* https://cloud.google.com/vision
* https://github.com/GoogleCloudPlatform/golang-samples/tree/master/vision/detect
* https://firebase.google.com/docs/storage/web/upload-files
* https://github.com/EddyVerbruggen/nativescript-plugin-firebase/blob/master/docs/FIRESTORE.md
* https://firebase.google.com/learn/codelabs/firestore-web?hl=en&continue=https%3A%2F%2Fcodelabs.developers.google.com%2F#5
* https://github.com/firebase
* https://github.com/googlearchive/friendlypix-web
* https://cloud.google.com/blog/products/application-development/serverless-in-action-building-a-simple-backend-with-cloud-firestore-and-cloud-functions
* https://flutter.dev/docs/get-started/test-drive?tab=vscode
* https://codelabs.developers.google.com/?cat=firebase
* https://github.com/GoogleCloudPlatform/golang-samples/tree/master/firestore/firestore_snippets
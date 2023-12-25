# Go RESTful API with Firebase

A simple Go CRUD app using the Gin web framework and connected to Firebase Realtime Datastore on the backend

## Instructions to deploy to Google Cloud Run

1. Install [Google Cloud CLI](https://cloud.google.com/sdk/docs/install)
2. Run `gcloud init`, set project corresponding to Firebase project name
3. Go into Google Cloud IAM dashboard and grant "Firebase Realtime Datastore Admin" and "Firebase Realtime Datastore Admin Service Agent" role to the Default compute service account
4. Run `gcloud run deploy` in the repo and follow steps (select region `us-central1`)

# Oura Ring API

[![Netlify Status](https://api.netlify.com/api/v1/badges/1c359a00-5e16-4d8b-bff4-28fce7c61966/deploy-status)](https://app.netlify.com/sites/thomasvn-oura/deploys)

Playing around with and visualizing Oura Ring data. Built using [Buildpacks](https://buildpacks.io/) and [Google Cloud Run Functions](https://github.com/GoogleCloudPlatform/functions-framework-go).

- <https://cloud.ouraring.com/dashboard>
- <https://cloud.ouraring.com/v2/docs>

## Demo

<https://oura.thomasvn.dev>

## Local Testing

```sh
pack build \
  --builder gcr.io/buildpacks/builder:v1 \
  --env GOOGLE_FUNCTION_SIGNATURE_TYPE=http \
  --env GOOGLE_FUNCTION_TARGET=Streaks \
  oura-streaks
```

```sh
docker run --rm -p 8080:8080 --env-file=.env oura-streaks
```

```sh
curl localhost:8080
```

## Deploy to Google Cloud Run Functions

```sh
gcloud functions deploy oura-streaks \
  --gen2 \
  --runtime=go123 \
  --region=us-west1 \
  --source=. \
  --entry-point=Streaks \
  --trigger-http \
  --allow-unauthenticated \
  --memory=128Mi \
  --cpu=.1 \
  --env-vars-file=.env.yaml

gcloud functions describe oura-streaks --region=us-west1
gcloud functions delete oura-streaks --region=us-west1
```

## References

- <https://cloud.google.com/functions/docs/create-deploy-http-go>
- <https://github.com/GoogleCloudPlatform/functions-framework-go>
- <https://cal-heatmap.com/>

<!-- 
IDEAS:
- Steps taken? Heatmap? Line graph?
- Calories Burned? Plot these together?
- Live streaming of biometric data? Webhook Subscription?
- Sedentary vs Low vs Medium vs High
- Bar graph. High + Medium + Low Activity Time.
- What's most important? Energy. What's most representative? HRV + Readiness.
-->

<!--
TODO:
- Stress versus Recovery. PlusMinus.
- Backend caching?
- Heatmap scrolling
-->

<!-- 
DONE (most recent first):
- Deploy to Netlify and SubDomain https://oura.thomasvn.dev
- Basic frontend for both these APIs
- Heatmap of Sleep/Readiness/Activity Scores (github-style)
- Deploy to Google Cloud Run Functions
- Local testing with `pack`
- Streak counter
  - Days above 75
  - Longest streak this year
- Successfully make an API call to Oura
-->

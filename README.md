# Oura Ring API

- <https://cloud.ouraring.com/dashboard>
- <https://cloud.ouraring.com/v2/docs>

```sh
export PAT=YOUR_OURA_PERSONAL_ACCESS_TOKEN
export LOCAL_ONLY=true

FUNCTION_TARGET=Streaks go run cmd/main.go
```

```sh
curl localhost:8080
```

<!-- 
IDEAS:
- Streak counter. Days above 75. "Longest streak this past year".
- Heatmap of Sleep/Readiness/Activity Scores (github-style)
- Live streaming of biometric data?
-->

<!--
TODO:
- Streak counter
  - Refactor for readability
  - Expose via an API
-->

<!-- 
DONE (most recent first):
- Streak counter
  - Days above 75
  - Longest streak this year
- Successfully make an API call to Oura
-->

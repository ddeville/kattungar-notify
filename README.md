# kattungar-notify

## Server credentials

### Google Calendar

- I created a new Google Cloud application and set up an OAuth 2.0 Client IDs credentials [here](https://console.cloud.google.com/apis/credentials?project=kattungar-notify).
- I used the [OAuth Playground](https://developers.google.com/oauthplayground) to generate a refresh token with the `https://www.googleapis.com/auth/calendar.events.readonly` and `https://www.googleapis.com/auth/calendar.readonly` scopes. Had to select "Use your own OAuth credentials" in the gear settings menu and input the ones from the app above.
- The `client_credentials.json` file was downloaded from `https://console.cloud.google.com/apis/credentials?project=kattungar-notify`.

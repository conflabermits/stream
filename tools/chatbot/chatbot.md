# How to chatbot

1. Set up this project, Go, etc.
2. Have a Twitch account.
3. Configure Twitch chat bot application in [Twitch Developer Console](https://dev.twitch.tv/console) and grab the Client ID and Secret.
4. Set up a `.creds` file like the `example.creds` file.
    * Fill in the twitchUsername, clientId, and clientSecret values, the rest can be left blank for now.
5. Get an OAuth token for the chatbot account.
    * Export the vars from the `.creds` file to your local env, go to the twitch-oauth-authorization-code-example directory, run main.go, open localhost:8080, and grab the token from the terminal after successful login.
    * Add the token value to the twitchToken line in the `.creds` file. It should looke like `twitchToken=oauth:YOUR_TOKEN_HERE`.
6. Use the token and client ID to get the broadcaster ID from the Twitch API.
    * Example command: `curl -sLkX GET "https://api.twitch.tv/helix/users?login=nightbot" -H "Authorization: Bearer ${twitchToken}" -H "Client-Id: ${clientId}" | jq .`
    * Add the value to the broadcaster_id line in the `.creds` file.
7. Fill in the name of the channel on the twitchChannel line in the `.creds` file.
8. Run the chatbot.
    * Export the vars from the `.creds` file to your local env, go to the chatgpt-chatbot-example directory, run main.go.

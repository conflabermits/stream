# Twitch Stream Chatbot

The purpose of this program is to act as a Twitch chatbot and respond to commmands on behalf of the streamer.

## Short term goals

* Combine the functionality of "twitch-token-example" and "chatgpt-chatbot-example" into a single program that seamlessly authenticates with Twitch to grab a fresh token, and to create a persistent connection to the Twitch channel, when given only a `.creds` file containing the necessary information (as seen in `example.creds`).

## Long term goals

* Merge the functionality of the existing chatbot and the donorbox overlay to have one single Go program that does all of the following:
  * Authenticate with Twitch and connect to a channel.
  * Respond to commands when posted in the chat.
  * Scrape the Donorbox site regularly to get the latest campaign info.
  * Host a local web server running the Donorbox HTML overlay.
  * Update the HTML (ideally with flashy graphics) when a new donation is received.
  * Respond to commands in the chat with the latest Donorbox campaign info when requested.
  * Send a message in the chat (ideally an announcement or something that triggers an overlay's alert sound) that notifies the chat when a new donation is received.

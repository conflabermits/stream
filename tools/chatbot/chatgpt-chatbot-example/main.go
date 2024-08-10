// Generated by ChatGPT and Gemini, mostly

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gempir/go-twitch-irc/v4"
)

func getEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		fmt.Printf("Error: Environment variable %s is not set\n", key)
		os.Exit(1)
	}
	return value
}

type PollPostData struct {
	// Define your data structure here
	BroadcasterId              string           `json:"broadcaster_id"`
	PollTitle                  string           `json:"title"`
	Choices                    []PollPostChoice `json:"choices"`
	ChannelPointsVotingEnabled bool             `json:"channel_points_voting_enabled"`
	ChannelPointsPerVote       int              `json:"channel_points_per_vote"`
	Duration                   int              `json:"duration"`
}

type PollPostChoice struct {
	Title string `json:"title"`
}

type PollGetResponse struct {
	Data []PollGetData `json:"data"`
}

type PollGetData struct {
	ID                         string          `json:"id"`
	BroadcasterID              string          `json:"broadcaster_id"`
	BroadcasterName            string          `json:"broadcaster_name"`
	BroadcasterLogin           string          `json:"broadcaster_login"`
	Title                      string          `json:"title"`
	Choices                    []PollGetChoice `json:"choices"`
	ChannelPointsVotingEnabled bool            `json:"channel_points_voting_enabled"`
	ChannelPointsPerVote       int             `json:"channel_points_per_vote"`
	Status                     string          `json:"status"`
	Duration                   int             `json:"duration"`
	StartedAt                  string          `json:"started_at"`
	EndedAt                    string          `json:"ended_at"`
}

type PollGetChoice struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Votes              int    `json:"votes"`
	ChannelPointsVotes int    `json:"channel_points_votes"`
}

func getPoll() (PollGetResponse, error) {
	// Define your URL and data
	url := "https://api.twitch.tv/helix/polls"

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return PollGetResponse{}, err
	}

	// Set the headers
	bearerToken := strings.Split(getEnvVar("twitchToken"), ":")[1]
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Client-Id", getEnvVar("clientId"))

	// Set the query parameters
	q := req.URL.Query()
	q.Add("broadcaster_id", getEnvVar("broadcaster_id"))
	q.Add("status", "ACTIVE")
	req.URL.RawQuery = q.Encode()

	// Send the request with a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return PollGetResponse{}, err
	}
	defer resp.Body.Close()

	// Optionally, handle the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return PollGetResponse{}, err
	}

	fmt.Println("Response status:", resp.StatusCode)

	// Parse the JSON data into a Data struct
	var jsonResponse PollGetResponse
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		fmt.Println("Error parsing JSON data", err)
		return PollGetResponse{}, err
	}

	return jsonResponse, nil
}

func getPollResults() string {
	jsonResponse, err := getPoll()
	if err != nil {
		errorMessage := "Error getting poll with getPoll()"
		fmt.Println(errorMessage, err)
		return errorMessage
	}

	var recentPoll string

	// Access the first element of the array
	if len(jsonResponse.Data) > 0 {
		firstValue := jsonResponse.Data[0]
		recentPoll += "Latest poll data! //// Poll title: \"" + firstValue.Title + "\" //// Choices:"
		for _, choice := range firstValue.Choices {
			recentPoll += " // \"" + choice.Title + "\": " + fmt.Sprint(choice.Votes)
		}
	} else {
		emptyMessage := "Array is empty"
		fmt.Println(emptyMessage)
		return emptyMessage
	}

	return recentPoll
}

func isPollActive() bool {
	jsonResponse, err := getPoll()
	if err != nil {
		fmt.Println("Error getting poll with getPoll()", err)
	}

	var isActive bool

	// Access the first element of the array
	if len(jsonResponse.Data) > 0 {
		firstValue := jsonResponse.Data[0]
		if firstValue.Status == "ACTIVE" {
			isActive = true
		} else {
			isActive = false
		}
	} else {
		fmt.Println("Array is empty")
		isActive = false
	}

	return isActive
}

func sendPoll(pollText string) string {
	log.Println("pollText: " + pollText)

	pollLength := len(strings.Split(pollText, "//"))
	log.Println("pollLength:", pollLength)
	if pollLength < 3 {
		errorMessage := "Not enough choices"
		fmt.Println(errorMessage)
		return errorMessage
	} else if pollLength > 6 {
		errorMessage := "Too many choices"
		fmt.Println(errorMessage)
		return errorMessage
	}

	var question string
	choices := []PollPostChoice{}

	for index, phrase := range strings.Split(pollText, "//") {
		if index == 0 {
			question = strings.TrimSpace(phrase)
			fmt.Println("Identified question as:", question)
		} else {
			fmt.Printf("Identified choice %s as: %s\n", fmt.Sprint(index), strings.TrimSpace(phrase))
			choices = append(choices, PollPostChoice{Title: strings.TrimSpace(phrase)})
		}
	}

	// Define your URL and data
	url := "https://api.twitch.tv/helix/polls"
	data := PollPostData{
		BroadcasterId:              getEnvVar("broadcaster_id"),
		PollTitle:                  question,
		Choices:                    choices,
		ChannelPointsVotingEnabled: true,
		ChannelPointsPerVote:       100,
		Duration:                   60,
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		errorMessage := "Error marshalling JSON"
		fmt.Println(errorMessage, err)
		return errorMessage
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		errorMessage := "Error creating request"
		fmt.Println(errorMessage, err)
		return errorMessage
	}

	// Set the headers
	bearerToken := strings.Split(getEnvVar("twitchToken"), ":")[1]
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Client-Id", getEnvVar("clientId"))

	// -H "Authorization: Bearer ${twitchToken}" -H "Client-Id: ${clientId}"

	// Send the request with a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorMessage := "Error sending request"
		fmt.Println(errorMessage, err)
		return errorMessage
	}
	defer resp.Body.Close()

	// Optionally, handle the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMessage := "Error reading response body"
		fmt.Println(errorMessage, err)
		return errorMessage
	}

	fmt.Println("Response status:", resp.StatusCode)
	fmt.Println("Response body:", string(body))

	if resp.StatusCode != 200 {
		return "Error creating poll"
	} else {
		return "Poll created successfully"
	}

}

// Borrowed code for better alphabetizer, case insensitive!
// https://programming-idioms.org/idiom/297/sort-a-list-of-strings-case-insensitively/5458/go

func lessCaseInsensitive(s, t string) bool {
	for {
		if len(t) == 0 {
			return false
		}
		if len(s) == 0 {
			return true
		}
		c, sizec := utf8.DecodeRuneInString(s)
		d, sized := utf8.DecodeRuneInString(t)

		lowerc := unicode.ToLower(c)
		lowerd := unicode.ToLower(d)

		if lowerc < lowerd {
			return true
		}
		if lowerc > lowerd {
			return false
		}

		s = s[sizec:]
		t = t[sized:]
	}
}

func alphabetize(message string) string {
	words := strings.Fields(message)
	sort.Slice(words, func(i, j int) bool { return lessCaseInsensitive(words[i], words[j]) })
	result := strings.Join(words, " ")
	return result
}

// Borrowed some code from Twilio's example for getting inspirational quotes
// https://www.twilio.com/blog/inspire-your-friends-using-go-twilio-messaging

type Response struct {
	Status     string
	StatusCode int
	Method     string
	Body       []byte
}

type ZenQuotes struct {
	Quote  string `json:"q"`
	Author string `json:"a"`
}

func getQuote() string {
	response, err := sendRequest()
	if err != nil {
		log.Fatal(err)
	}

	var message ZenQuotes
	msg := []ZenQuotes{{
		message.Author,
		message.Quote,
	}}

	err = json.Unmarshal(response.Body, &msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg[0].Quote
}

func sendRequest() (*Response, error) {
	r := &Response{}

	httpClient := &http.Client{Timeout: 20 * time.Second}
	zenQuotesUrl := "https://zenquotes.io/api/random"

	req, err := http.NewRequest(http.MethodGet, zenQuotesUrl, nil)
	if err != nil {
		return nil, err
	}

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	r.Status = response.Status
	r.StatusCode = response.StatusCode
	r.Body = body

	return r, nil
}

func main() {
	// Set your Twitch username, OAuth token, and channel to join
	//username := "your_twitch_username"
	//token := "your_oauth_token"
	//channel := "channel_to_join"
	username := getEnvVar("twitchUsername") //e.g., "conflabermits"
	token := getEnvVar("twitchToken")       //e.g., "oauth:<token>"
	channel := getEnvVar("twitchChannel")   //e.g., "conflabermits"

	// Create a new Twitch client
	client := twitch.NewClient(username, token)

	// Register a callback for when the bot receives a message
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		// Print the message to the console
		fmt.Printf("[%s] %s: %s\n", message.Channel, message.User.DisplayName, message.Message)

		// You can add your own logic here to respond to messages
		// For example, you can check for specific commands and reply accordingly
		// TODO: Make a function to parse messages using common conditions. Examples:
		//   * startsWith(message.Message, "!string ")
		//   * equals(message.Message, "!string")
		//   * hasArg(message.Message, [eq,lt,gt], int)
		//   * fromUser(message.Message, "user")
		//   * fromRole(message.Message, [bc,mod,vip,sub,fol])
		if message.Message == "!hello" || message.Message == "!hellobot" {
			log.Println("Detected !hello message")
			client.Say(message.Channel, "Hello, "+message.User.DisplayName+"!")
		}
		if message.Message == "!bye" || message.Message == "!byebot" {
			log.Println("Detected !bye message")
			client.Say(message.Channel, "Goodbye, "+message.User.DisplayName+"! I'll miss you!")
		}
		if strings.HasPrefix(message.Message, "!abc ") || strings.HasPrefix(message.Message, "!alpha ") {
			log.Println("Detected !abc message")
			commandText := strings.TrimPrefix(message.Message, "!abc ")
			//TODO: Ensure there is a message after the !abc command to be alphabetized
			client.Say(message.Channel, alphabetize(commandText))
		}
		// Command ideas:
		// !randomize - Randomize the words from the given message.
		// !lore - Print a random line from a text file containing deep conflabermits lore.
		if message.Message == "!quote" || message.Message == "!randomquote" {
			log.Println("Detected !quote message")
			client.Say(message.Channel, "Random quote -- "+getQuote()+".. in bed.")
		}
		if message.Message == "!poll" || message.Message == "!getPoll" || message.Message == "!getPollResults" {
			log.Println("Detected !getPoll message")
			client.Say(message.Channel, getPollResults())
		}
		if strings.HasPrefix(message.Message, "!poll ") {
			log.Println("Detected !poll message")
			if isPollActive() {
				client.Say(message.Channel, "Sorry, a poll is currently active, try again when it's done.")
			} else if message.User.DisplayName != "conflabermits" {
				client.Say(message.Channel, "Sorry, only accepting polls from conflabermits right now!")
			} else {
				client.Say(message.Channel, "Attempting to create a poll for @"+message.User.DisplayName+"...")
				pollText := strings.TrimPrefix(message.Message, "!poll ")
				client.Say(message.Channel, sendPoll(pollText))
			}
		}
	})

	client.OnConnect(func() {
		client.Say(channel, "conflabermits chatbot started -- Let's GOOOOOO! confla3Nyan confla3Shake confla3Party confla3Spin")
	})

	// Join the specified channel
	client.Join(channel)

	// Connect to Twitch IRC
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for a signal to gracefully shut down the bot
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Disconnect from Twitch IRC on shutdown
	client.Disconnect()
}

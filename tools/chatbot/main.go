package main

import (
	"context"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

/* GH CoPilot's suggestion for combining the two programs

func getOAuthToken() string {
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/login", handleTwitchLogin)
    http.HandleFunc("/oauth2", handleTwitchCallback)

    go http.ListenAndServe(":8080", nil)

    // Wait for auth to complete
    token := <-authComplete
    return token.AccessToken
}

// Rest of the first program's functions...

// Second program's functions...

func main() {
    // Set your Twitch username, OAuth token, and channel to join
    username := getEnvVar("twitchUsername") //e.g., "conflabermits"
    token := getOAuthToken()                 // Get OAuth token
    channel := getEnvVar("twitchChannel")    //e.g., "conflabermits"

    // Rest of the second program's main function...
}
*/

// Twitch OAuth2 Code

const (
	stateCallbackKey = "oauth-state-callback"
	oauthSessionName = "oauth-session"
	oauthTokenKey    = "oauth-token"
)

var (
	//clientID = "<CLIENT_ID>"
	//clientSecret = "<CLIENT_SECRET>"
	// Consider storing the secret in an environment variable or a dedicated storage system.
	clientID     = getEnvVar("clientId")
	clientSecret = getEnvVar("clientSecret")
	scopes       = []string{"chat:read", "chat:edit"}
	redirectURL  = "http://localhost:8080/redirect"
	oauth2Config *oauth2.Config
	cookieSecret = []byte("I don't think this is a good secret either")
	cookieStore  = sessions.NewCookieStore(cookieSecret)
	oAuthToken   = ""
)

func getEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		fmt.Printf("Error: Environment variable %s is not set\n", key)
		os.Exit(1)
	}
	return value
}

// HandleRoot is a Handler that shows a login button. In production, if the frontend is served / generated
// by Go, it should use html/template to prevent XSS attacks.
func HandleRoot(w http.ResponseWriter, r *http.Request) (err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<html><body><a href="/login">Login using Twitch</a></body></html>`))

	return
}

// HandleLogin is a Handler that redirects the user to Twitch for login, and provides the 'state'
// parameter which protects against login CSRF.
func HandleLogin(w http.ResponseWriter, r *http.Request) (err error) {
	session, err := cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	var tokenBytes [255]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		return AnnotateError(err, "Couldn't generate a session!", http.StatusInternalServerError)
	}

	state := hex.EncodeToString(tokenBytes[:])

	session.AddFlash(state, stateCallbackKey)

	if err = session.Save(r, w); err != nil {
		return
	}

	http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusTemporaryRedirect)

	return
}

// HandleOauth2Callback is a Handler for oauth's 'redirect_uri' endpoint;
// it validates the state token and retrieves an OAuth token from the request parameters.
func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (err error) {
	session, err := cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	// ensure we flush the csrf challenge even if the request is ultimately unsuccessful
	defer func() {
		if err := session.Save(r, w); err != nil {
			log.Printf("error saving session: %s", err)
		}
	}()

	switch stateChallenge, state := session.Flashes(stateCallbackKey), r.FormValue("state"); {
	case state == "", len(stateChallenge) < 1:
		err = errors.New("missing state challenge")
	case state != stateChallenge[0]:
		err = fmt.Errorf("invalid oauth state, expected '%s', got '%s'\n", state, stateChallenge[0])
	}

	if err != nil {
		return AnnotateError(
			err,
			"Couldn't verify your confirmation, please try again.",
			http.StatusBadRequest,
		)
	}

	token, err := oauth2Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return
	}

	// add the oauth token to session
	session.Values[oauthTokenKey] = token

	fmt.Printf("Access token: %s\n", token.AccessToken)
	//fmt.Printf("Access token got\n")

	// store the token in a variable for use by the chatbot
	oAuthToken = token.AccessToken

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	return
}

// HumanReadableError represents error information
// that can be fed back to a human user.
//
// This prevents internal state that might be sensitive
// being leaked to the outside world.
type HumanReadableError interface {
	HumanError() string
	HTTPCode() int
}

// HumanReadableWrapper implements HumanReadableError
type HumanReadableWrapper struct {
	ToHuman string
	Code    int
	error
}

func (h HumanReadableWrapper) HumanError() string { return h.ToHuman }
func (h HumanReadableWrapper) HTTPCode() int      { return h.Code }

// AnnotateError wraps an error with a message that is intended for a human end-user to read,
// plus an associated HTTP error code.
func AnnotateError(err error, annotation string, code int) error {
	if err == nil {
		return nil
	}
	return HumanReadableWrapper{ToHuman: annotation, error: err}
}

type Handler func(http.ResponseWriter, *http.Request) error

func getOAuthToken() {

	/*
		https://github.com/twitchdev/authentication-go-sample/blob/main/oauth-authorization-code/main.go
		Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
		Licensed under the Apache License, Version 2.0 (the "License").
	*/

	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})

	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirectURL,
	}

	var middleware = func(h Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) (err error) {
			// parse POST body, limit request size
			if err = r.ParseForm(); err != nil {
				return AnnotateError(err, "Something went wrong! Please try again.", http.StatusBadRequest)
			}

			return h(w, r)
		}
	}

	// errorHandling is a middleware that centralises error handling.
	// this prevents a lot of duplication and prevents issues where a missing
	// return causes an error to be printed, but functionality to otherwise continue
	// see https://blog.golang.org/error-handling-and-go
	var errorHandling = func(handler func(w http.ResponseWriter, r *http.Request) error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := handler(w, r); err != nil {
				var errorString string = "Something went wrong! Please try again."
				var errorCode int = 500

				if v, ok := err.(HumanReadableError); ok {
					errorString, errorCode = v.HumanError(), v.HTTPCode()
				}

				log.Println(err)
				w.Write([]byte(errorString))
				w.WriteHeader(errorCode)
				return
			}
		})
	}

	var handleFunc = func(path string, handler Handler) {
		http.Handle(path, errorHandling(middleware(handler)))
	}

	handleFunc("/", HandleRoot)
	handleFunc("/login", HandleLogin)
	handleFunc("/redirect", HandleOAuth2Callback)

	fmt.Println("Started running on http://localhost:8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}

// Twitch Chatbot Code

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

func startChatbot() {
	// Set your Twitch username, OAuth token, and channel to join
	//username := "your_twitch_username"
	//token := "your_oauth_token"
	//channel := "channel_to_join"
	username := getEnvVar("twitchUsername") //e.g., "conflabermits"
	channel := getEnvVar("twitchChannel")   //e.g., "conflabermits"
	token := oAuthToken                     //e.g., "oauth:<token>"

	// Create a new Twitch client
	client := twitch.NewClient(username, token)

	// Register a callback for when the bot receives a message
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		// Print the message to the console
		fmt.Printf("[%s] %s: %s\n", message.Channel, message.User.DisplayName, message.Message)

		// You can add your own logic here to respond to messages
		// For example, you can check for specific commands and reply accordingly
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
			client.Say(message.Channel, alphabetize(commandText))
		}
		// Command ideas:
		// !randomize - Randomize the words from the message
		if message.Message == "!quote" || message.Message == "!randomquote" {
			log.Println("Detected !quote message")
			client.Say(message.Channel, "Random quote -- "+getQuote()+".. in bed.")
		}
	})

	client.OnConnect(func() { client.Say(channel, "Let's GOOOOOO!") })

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

func main() {
	go getOAuthToken()
	startChatbot()
}

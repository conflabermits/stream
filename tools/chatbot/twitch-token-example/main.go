// https://github.com/twitchdev/authentication-go-sample/blob/main/oauth-client-credentials/main.go

/*
Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may not use this file except in compliance with the License. A copy of the License is located at

	http://aws.amazon.com/apache2.0/

	or in the "license" file accompanying this file. This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

*/

package main

import (
	"context"
	"fmt"
	"os"
	"log"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

func getEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		fmt.Printf("Error: Environment variable %s is not set\n", key)
		os.Exit(1)
	}
	return value
}

var (
	//clientID = "<CLIENT_ID>"
	clientID = getEnvVar("clientId")
	// Consider storing the secret in an environment variable or a dedicated storage system.
	//clientSecret = "<CLIENT_SECRET>"
	clientSecret = getEnvVar("clientSecret")
	oauth2Config *clientcredentials.Config
)

func main() {
	oauth2Config = &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	token, err := oauth2Config.Token(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", token.AccessToken)
}


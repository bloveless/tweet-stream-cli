package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"tweet-stream-go/pkg/twitter"
)

func main() {
	c, err := loadConfig()
	if err != nil {
		panic(err)
	}

	// requestToken(c)
	// postStatus(c)
	streamTweets(c)
}

type config struct {
	OAuthConsumerKey       string `json:"oauth_consumer_key"`
	OAuthConsumerSecret    string `json:"oauth_consumer_secret"`
	OAuthAccessToken       string `json:"oauth_access_token"`
	OAuthAccessTokenSecret string `json:"oauth_access_token_secret"`
}

func getConfigPath() (string, error) {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = os.Getenv("HOME") + "/.config"
	}

	configDirectory := xdgConfigHome + "/tweet-stream-cli"

	err := os.MkdirAll(configDirectory, 0755)
	if err != nil {
		return "", err
	}

	return configDirectory + "/config.json", nil
}

func loadConfig() (config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return config{}, err
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config{}, err
	}

	c := config{}
	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		return config{}, err
	}

	return c, nil
}

func saveConfig(c config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func getTwitterClient(c config) twitter.Client {
	return twitter.Client{
		OAuthConsumerKey:       c.OAuthConsumerKey,
		OAuthConsumerSecret:    c.OAuthConsumerSecret,
		OAuthAccessToken:       c.OAuthAccessToken,
		OAuthAccessTokenSecret: c.OAuthAccessTokenSecret,
		// HTTPClient: http.Client{
		// 	Timeout: 1 * time.Minute,
		// },
	}
}

func requestToken(c config) {
	tc := getTwitterClient(c)

	input := twitter.OAuthRequestTokenInput{
		OAuthCallback: "http://127.0.0.1:3000/oauth_response",
	}

	output, err := tc.OAuthRequestTokenGet(input)
	if err != nil {
		panic(err)
	}

	fmt.Println("Please authenticate by visiting this link https://api.twitter.com/oauth/authorize?oauth_token=" + output.OAuthToken)

	type oAuthResponse struct {
		OAuthToken    string
		OAuthVerifier string
	}

	done := make(chan oAuthResponse, 1)
	shutdown := make(chan bool, 1)

	router := http.NewServeMux()
	router.HandleFunc("/oauth_response", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Access has been successfully recorded. You can return to the comfort of your terminal"))

		or := oAuthResponse{
			OAuthToken:    r.URL.Query().Get("oauth_token"),
			OAuthVerifier: r.URL.Query().Get("oauth_verifier"),
		}

		done <- or
	})

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	listenAddr := ":3000"
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	or := oAuthResponse{}
	go func() {
		or = <-done
		logger.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(shutdown)
	}()

	logger.Println("Server is waiting for auth token at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-shutdown
	logger.Println("Server stopped")
	logger.Println(or)

	accessTokenInput := twitter.OAuthAccessTokenInput{
		OAuthToken:    or.OAuthToken,
		OAuthVerifier: or.OAuthVerifier,
	}
	accessTokenOutput, err := tc.OAuthAccessTokenGet(accessTokenInput)
	if err != nil {
		panic(err)
	}

	fmt.Println(accessTokenOutput)
	updatedConfig := config{
		OAuthConsumerKey:       c.OAuthConsumerKey,
		OAuthConsumerSecret:    c.OAuthConsumerSecret,
		OAuthAccessToken:       accessTokenOutput.OAuthToken,
		OAuthAccessTokenSecret: accessTokenOutput.OAuthTokenSecret,
	}

	err = saveConfig(updatedConfig)
	if err != nil {
		panic(err)
	}
}

func postStatus(c config) {
	tc := getTwitterClient(c)

	input := twitter.StatusesUpdateInput{
		Status:            "Hello Ladies + Gentlemen, a signed OAuth request!",
		InReplyToStatusID: 0,
	}

	output, err := tc.StatusesUpdatePost(input) // GetSignedRequest(http.MethodPost, uri, params)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%+v\n", output)
}

func streamTweets(c config) {
	tc := getTwitterClient(c)

	input := twitter.StatusesFilterInput{
		Track: "kubernetes",
	}

	output, err := tc.StatusesFilterPost(input)
	if err != nil {
		panic(err)
	}

	count := 0
	for {
		tweet := twitter.StatusesFilterOutput{}
		err := json.NewDecoder(output.Body).Decode(&tweet)
		if err == io.EOF {
			fmt.Println("End of file")
		}

		if err != nil {
			panic(err)
		}

		fmt.Println(count)
		fmt.Printf("%#v\n\n\n", tweet)

		count++
		if count > 10 {
			return
		}
	}
}

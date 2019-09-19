package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kurrik/oauth1a"

	"github.com/kurrik/twittergo"

	"github.com/smugcloud/twitter-cleanup/twitter"
)

var handle, from string
var start int

var defaultSearch = `https://api.twitter.com/1.1/tweets/search/fullarchive/dev.json?query=from:`
var defaultDelete = `https://api.twitter.com/1.1/statuses/destroy/`

func main() {
	lambda.Start(DeleteTweets)
}

// DeleteTweets will handle the Lambda invocation
func DeleteTweets(ctx context.Context) error {
	handle = os.Getenv("HANDLE")
	start, _ = strconv.Atoi(os.Getenv("START"))
	from = os.Getenv("FROM")
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	tokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || tokenSecret == "" {
		log.Print("Environment variables for Twitter authentication are required.\n\n")
		os.Exit(1)
	}

	oauthConfig := &oauth1a.ClientConfig{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
	}
	user := &oauth1a.UserConfig{
		AccessTokenKey:    accessToken,
		AccessTokenSecret: tokenSecret,
	}

	tgoClient := twittergo.NewClient(oauthConfig, user)
	client := twitter.Client{
		Tgo:       tgoClient,
		SearchURL: defaultSearch,
		DeleteURL: defaultDelete,
		APIRequest: twitter.APIRequest{
			Handle:     handle,
			MonthsBack: start,
			From:       from,
		},
		DeleteIDS: make(chan uint64, 20),
	}
	var wg sync.WaitGroup
	wg.Add(1)
	client.ProcessTweets()
	wg.Wait()

	return nil
}

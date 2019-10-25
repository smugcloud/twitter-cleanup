package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"github.com/smugcloud/twitter-cleanup/twitter"
)

var handle, from string
var start int

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
	api := twitter.APIRequest{
		Handle:     handle,
		MonthsBack: start,
		From:       from,
	}

	client := twitter.NewClient(tgoClient, api)

	client.ProcessTweets()

	log.Print("Waitgroup released, returning no errors.")
	return nil
}

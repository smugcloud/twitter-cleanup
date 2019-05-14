package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kurrik/oauth1a"

	"github.com/kurrik/twittergo"

	"github.com/smugcloud/twitter-cleanup/server"
	"github.com/smugcloud/twitter-cleanup/twitter"
)

var handle, from string
var start, period int

var defaultSearch = `https://api.twitter.com/1.1/tweets/search/fullarchive/dev.json?query=from:`
var defaultDelete = `https://api.twitter.com/1.1/statuses/destroy/`

func main() {
	flag.StringVar(&handle, "handle", "smugcloud", "Twitter username to search.")
	flag.IntVar(&start, "start", 12, "The number of previous months to preserve on Twitter (anything before will be deleted).")
	flag.StringVar(&from, "from", "20100101", "Starting date to look back to, in Twitter format (YYYYMMDD)")
	flag.IntVar(&period, "period", 1, "The frequency with which to check Twitter (in months)")

	// Custom Usage function to show the environment variables needed.
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", os.Args[0])
		fmt.Printf("Environment variables for Twitter authentication:\n")
		fmt.Printf("CONSUMER_SECRET\nCONSUMER_KEY\nACCESS_TOKEN\nACCESS_TOKEN_SECRET\n\n")
		fmt.Println("Command line flags:")
		flag.PrintDefaults()

	}

	flag.Parse()
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	tokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || tokenSecret == "" {
		fmt.Printf("Environment variables for Twitter authentication are required.\n\n")
		flag.Usage()
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
			Period:     period,
		},
		DeleteIDS: make(chan uint64, 20),
	}
	client.ProcessTweets()

	s := server.Server{}
	s.HttpServer()
}

package main

import (
	"flag"
	"os"

	"github.com/kurrik/oauth1a"

	"github.com/kurrik/twittergo"

	"github.com/smugcloud/twitter-cleanup/server"
	"github.com/smugcloud/twitter-cleanup/twitter"
)

var handle, from string
var start, period int

func main() {
	flag.Parse()
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	tokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")

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
		Tgo: tgoClient,
		APIRequest: twitter.APIRequest{
			Handle:     handle,
			MonthsBack: start,
			Period:     period,
		},
		DeleteIDS: make(chan uint64, 20),
	}
	client.ProcessTweets()

	s := server.Server{}
	s.HttpServer()
}

func init() {
	flag.StringVar(&handle, "handle", "smugcloud", "Twitter username to search.")
	flag.IntVar(&start, "start", 12, "The number of previous months to preserve on Twitter (anything before will be deleted).")
	flag.StringVar(&from, "from", "20100101", "Starting date to look back to, in Twitter format (YYYYMMDD)")

	flag.IntVar(&period, "period", 1, "The frequency with which to check Twitter (in months)")

}

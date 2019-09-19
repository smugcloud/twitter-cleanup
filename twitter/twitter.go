package twitter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/smugcloud/twitter-cleanup/util"

	"github.com/kurrik/twittergo"
)

// ITwitter is an interface for mocking
type ITwitter interface {
	deleteTweets()
}

//Auth holds the bearer (acess) token we get back from Twitter
type Auth struct {
	AccessToken string `json:"access_token"`
}

//Response is the container to hold the Tweet ID's for processing
type Response struct {
	Results []struct {
		ID uint64 `json:"id"`
	} `json:"results"`
	Next string `json:"next"`
}

//Client is the shared struct to auth to Twitter
type Client struct {
	Tgo *twittergo.Client
	APIRequest
	DeleteIDS chan uint64
	SearchURL string
	DeleteURL string
}

//APIRequest holds the values we want to send in the request with multiple calls
type APIRequest struct {
	Handle     string
	To         string
	From       string
	Next       string
	MonthsBack int
	Period     int
	count      int
}

//NewITwitter provides an Interface to use for mocking/testing
func NewITwitter(c *Client) ITwitter {
	return c
}

//Cleanup is the ticker which triggers a new run of the ProcessTweets
func Cleanup(client Client) {
	// ~Hours in a month: 24 * 30 = 720
	t := time.NewTicker(time.Duration(client.Period*720) * time.Hour)
	for {
		select {
		case <-t.C:
			client.count = 0
			client.getAllTweets(&client.APIRequest)
			log.Printf("%v tweets processed.\n", client.count)
		}
	}
}

//ProcessTweets serves as the entry for program startup, as well as periodic cleanup.
func (c *Client) ProcessTweets() {
	c.APIRequest.To = util.GetToDate(c.APIRequest.MonthsBack, time.Now())

	// Watch the channel in a goroutine

	go c.deleteTweets()
	c.getAllTweets(&c.APIRequest)

}

func (c *Client) getAllTweets(options *APIRequest) {
	//While there's a `next` value in the response, follow the next, but also grab
	//the Tweet ID for deletion

	u := c.SearchURL + options.Handle + `&fromDate=` + options.From + `0000&toDate=` + options.To + `&maxResults=10`
	if options.Next != "" {
		u = u + "&next=" + options.Next
	}
	req, _ := http.NewRequest("GET", u, nil)

	resp, err := c.Tgo.SendRequest(req)
	if err != nil {
		log.Fatalf("error in querying tweets: %v\n", err)
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)

		log.Printf("Full response on non 200: %v\n", resp)
		log.Printf("Response body on non 200: %v\n", string(body))
	}
	r := Response{}
	json.NewDecoder(resp.Body).Decode(&r)

	// Put all of the latest results on the channel
	for _, v := range r.Results {
		c.count++
		c.DeleteIDS <- v.ID

	}
	if r.Next != "" {
		c.APIRequest.Next = r.Next
		c.getAllTweets(&c.APIRequest)
	}
	close(c.DeleteIDS)
}

// Function to watch the channel, and delete all the tweets that are on it.
func (c *Client) deleteTweets() {
	var count int
	for id := range c.DeleteIDS {
		count++
		// Manually appending the slash for now
		u := c.DeleteURL + strconv.FormatUint(id, 10) + ".json"
		req, _ := http.NewRequest("POST", u, nil)

		resp, err := c.Tgo.SendRequest(req)
		if err != nil {
			log.Fatalf("error in delete request: %v\n", err)
		}
		if resp.StatusCode != 200 {
			log.Printf("Received a non-200 value in delete: %v\n", resp.StatusCode)
			body, _ := ioutil.ReadAll(resp.Body)

			log.Printf("Full response in delete: %v\n", resp)
			log.Printf("Response body in delete: %v\n", string(body))
			continue
		}
	}
	log.Printf("Deleted %v total tweets.", count)

}

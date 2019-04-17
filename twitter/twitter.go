package twitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/smugcloud/twitter-cleanup/util"

	"github.com/kurrik/twittergo"
)

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
}

//APIRequest holds the values we want to send in the request with multiple calls
type APIRequest struct {
	Handle     string
	To         string
	Next       string
	MonthsBack int
	Period     int
}

//Cleanup is the ticker which triggers a new run of the ProcessTweets
func Cleanup(client Client) {
	// ~Hours in a month: 24 * 30 = 720
	t := time.NewTicker(time.Duration(client.Period*720) * time.Hour)
	for {
		select {
		case <-t.C:
			// client.getBearerToken()
			client.getAllTweets(&client.APIRequest)
		}
	}
}

//ProcessTweets serves as the entry for program startup, as well as periodic cleanup.
func (c *Client) ProcessTweets() {
	c.APIRequest.To = util.GetToDate(c.APIRequest.MonthsBack, time.Now())

	c.getAllTweets(&c.APIRequest)
	c.deleteTweets()
}

func (c *Client) getAllTweets(options *APIRequest) {
	//While there's a `next` value in the response, follow the next, but also grab
	//the Tweet ID for deletion

	u := `https://api.twitter.com/1.1/tweets/search/fullarchive/dev.json?query=@` + options.Handle + `&toDate=` + options.To + `&maxResults=10`
	if options.Next != "" {
		u = u + "&next=" + options.Next
	}
	// log.Printf("Url to use: %v\n", u)
	req, _ := http.NewRequest("GET", u, nil)

	resp, err := c.Tgo.SendRequest(req)
	if err != nil {
		log.Fatalf("error in Do: %v\n", err)
	}
	log.Printf("getTweets: %+v", resp)
	r := Response{}
	json.NewDecoder(resp.Body).Decode(&r)

	for _, v := range r.Results {
		c.DeleteIDS <- v.ID

	}
	log.Printf("Results: %+v", r)
	if r.Next != "" {
		c.APIRequest.Next = r.Next
		c.getAllTweets(&c.APIRequest)
	}
}

func (c *Client) deleteTweets() {
	fmt.Println("Inside delete")
	for {
		select {
		case id := <-c.DeleteIDS:
			fmt.Printf("Would delete %v\n", id)
			// u := "https://api.twitter.com/1.1/statuses/destroy/" + string(id) + ".json"

			// req, _ := http.NewRequest("POST", u, nil)

			// resp, err := c.Tgo.SendRequest(req)
			// if err != nil {
			// 	log.Fatalf("error in Do: %v\n", err)
			// }
			// if resp.StatusCode != 200 {

			// }
		}
	}

}

// func (c *Client) getBearerToken() {
// 	//Call the API and get the AccessToken
// 	u := "https://api.twitter.com/oauth2/token"

// 	client := http.Client{}
// 	// read := strings.NewReader("grant_type=client_credentials")
// 	data := url.Values{}
// 	data.Set("grant_type", "client_credentials")

// 	req, _ := http.NewRequest("POST", u, strings.NewReader(data.Encode()))
// 	err := req.ParseForm()
// 	if err != nil {
// 		log.Fatalf("error parsing form: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

// 	req.SetBasicAuth(c.Username, c.Password)
// 	resp, err := client.Do(req)
// 	auth := Auth{}
// 	json.NewDecoder(resp.Body).Decode(&auth)
// 	c.Bearer = auth.AccessToken

// }

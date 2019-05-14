package twitter

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"github.com/smugcloud/twitter-cleanup/util"
)

func NewTGo() *twittergo.Client {
	oauthConfig := &oauth1a.ClientConfig{
		ConsumerKey:    "consumerKey",
		ConsumerSecret: "consumerSecret",
	}
	user := &oauth1a.UserConfig{
		AccessTokenKey:    "accessToken",
		AccessTokenSecret: "tokenSecret",
	}

	tgoClient := twittergo.NewClient(oauthConfig, user)
	return tgoClient
}

func newTestClient() *Client {
	return &Client{
		DeleteIDS: make(chan uint64, 1),
	}
}
func TestDeleteTweets(t *testing.T) {
	id := uint64(1058408022936977409)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received: %v\n", r.URL.Path)

		period := strings.Index(r.URL.Path, ".")
		if strconv.FormatUint(id, 10) != r.URL.Path[1:period] {
			t.Errorf("Incorrect path requested to Twitter: %v", r.URL.Path)
		}
	})

	ts := httptest.NewServer(handler)
	log.Printf("Server started at %v\n", ts.URL)

	defer ts.Close()
	// Fill the auth values
	tgo := NewTGo()

	// Get a client we can work with
	client := newTestClient()
	// Fill the client values
	client.Tgo = tgo
	client.SearchURL = util.URLParse(ts.URL).String()
	client.DeleteURL = util.URLParse(ts.URL).String()

	iTwitter := NewITwitter(client)

	go iTwitter.deleteTweets()
	//Send an ID on the channel
	client.DeleteIDS <- id
	time.Sleep(1 * time.Second)

}

func TestGetTweets(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received: %v\n", r.URL.Path)

		w.Write([]byte(apiResponse))
	})

	ts := httptest.NewServer(handler)
	log.Printf("Server started at %v\n", ts.URL)

	defer ts.Close()
	// Fill the auth values
	tgo := NewTGo()

	// Get a client we can work with
	client := newTestClient()
	// Fill the client values
	client.Tgo = tgo
	client.SearchURL = util.URLParse(ts.URL).String()
	client.DeleteURL = util.URLParse(ts.URL).String()
	options := &APIRequest{}
	client.getAllTweets(options)

	if 1058408022936977409 != <-client.DeleteIDS {
		t.Error("Error in getAllTweets.")
	}
}

// Sample response from Twitter
const apiResponse = `{
	"results": [
		{
			"created_at": "Fri Nov 02 17:18:31 +0000 2018",
			"id": 1058408022936977409,
			"id_str": "1058408022936977409",
			"text": "RT @harmophone: \"The innovative crowdsourcing that the Tagboard, Twitter and TEGNA collaboration enables is surfacing locally relevant conv…",
			"source": "<a href=\"http:\/\/twitter.com\" rel=\"nofollow\">Twitter Web Client<\/a>",
			"truncated": false,
			"in_reply_to_status_id": null,
			"in_reply_to_status_id_str": null,
			"in_reply_to_user_id": null,
			"in_reply_to_user_id_str": null,
			"in_reply_to_screen_name": null,
			"user": {
				"id": 2244994945,
				"id_str": "2244994945",
				"name": "Twitter Dev",
				"screen_name": "TwitterDev",
				"location": "Internet",
				"url": "https:\/\/developer.twitter.com\/",
				"description": "Your official source for Twitter Platform news, updates & events. Need technical help? Visit https:\/\/twittercommunity.com\/ ⌨️ #TapIntoTwitter",
				"translator_type": "regular",
				"protected": false,
				"verified": true,
				"followers_count": 503828,
				"friends_count": 1477,
				"listed_count": 1437,
				"favourites_count": 2199,
				"statuses_count": 3380,
				"created_at": "Sat Dec 14 04:35:55 +0000 2013",
				"utc_offset": null,
				"time_zone": null,
				"geo_enabled": true,
				"lang": "en",
				"contributors_enabled": false,
				"is_translator": false,
				"profile_background_color": "FFFFFF",
				"profile_background_image_url": "http:\/\/abs.twimg.com\/images\/themes\/theme1\/bg.png",
				"profile_background_image_url_https": "https:\/\/abs.twimg.com\/images\/themes\/theme1\/bg.png",
				"profile_background_tile": false,
				"profile_link_color": "0084B4",
				"profile_sidebar_border_color": "FFFFFF",
				"profile_sidebar_fill_color": "DDEEF6",
				"profile_text_color": "333333",
				"profile_use_background_image": false,
				"profile_image_url": "http:\/\/pbs.twimg.com\/profile_images\/880136122604507136\/xHrnqf1T_normal.jpg",
				"profile_image_url_https": "https:\/\/pbs.twimg.com\/profile_images\/880136122604507136\/xHrnqf1T_normal.jpg",
				"profile_banner_url": "https:\/\/pbs.twimg.com\/profile_banners\/2244994945\/1498675817",
				"default_profile": false,
				"default_profile_image": false,
				"following": null,
				"follow_request_sent": null,
				"notifications": null
			},
			"geo": null,
			"coordinates": null,
			"place": null,
			"contributors": null,
			"retweeted_status": {
				"created_at": "Tue Oct 30 21:30:25 +0000 2018",
				"id": 1057384253116289025,
				"id_str": "1057384253116289025",
				"text": "\"The innovative crowdsourcing that the Tagboard, Twitter and TEGNA collaboration enables is surfacing locally relev… https:\/\/t.co\/w46U5TRTzQ",
				"source": "<a href=\"http:\/\/twitter.com\" rel=\"nofollow\">Twitter Web Client<\/a>",
				"truncated": true,
				"in_reply_to_status_id": null,
				"in_reply_to_status_id_str": null,
				"in_reply_to_user_id": null,
				"in_reply_to_user_id_str": null,
				"in_reply_to_screen_name": null,
				"user": {
					"id": 175187944,
					"id_str": "175187944",
					"name": "Tyler Singletary",
					"screen_name": "harmophone",
					"location": "San Francisco, CA",
					"url": "http:\/\/medium.com\/@harmophone",
					"description": "SVP Product at @Tagboard. Did some Data, biz, and product @Klout & for @LithiumTech; @BBI board member; @Insightpool advisor. World's worst whiteboarder.",
					"translator_type": "none",
					"protected": false,
					"verified": false,
					"followers_count": 1982,
					"friends_count": 1877,
					"listed_count": 245,
					"favourites_count": 23743,
					"statuses_count": 12708,
					"created_at": "Thu Aug 05 22:59:29 +0000 2010",
					"utc_offset": null,
					"time_zone": null,
					"geo_enabled": false,
					"lang": "en",
					"contributors_enabled": false,
					"is_translator": false,
					"profile_background_color": "C42329",
					"profile_background_image_url": "http:\/\/abs.twimg.com\/images\/themes\/theme1\/bg.png",
					"profile_background_image_url_https": "https:\/\/abs.twimg.com\/images\/themes\/theme1\/bg.png",
					"profile_background_tile": true,
					"profile_link_color": "E67D69",
					"profile_sidebar_border_color": "000000",
					"profile_sidebar_fill_color": "DDEEF6",
					"profile_text_color": "333333",
					"profile_use_background_image": true,
					"profile_image_url": "http:\/\/pbs.twimg.com\/profile_images\/719985428632240128\/WYFHcK-m_normal.jpg",
					"profile_image_url_https": "https:\/\/pbs.twimg.com\/profile_images\/719985428632240128\/WYFHcK-m_normal.jpg",
					"profile_banner_url": "https:\/\/pbs.twimg.com\/profile_banners\/175187944\/1398653841",
					"default_profile": false,
					"default_profile_image": false,
					"following": null,
					"follow_request_sent": null,
					"notifications": null
				},
				"geo": null,
				"coordinates": null,
				"place": null,
				"contributors": null,
				"is_quote_status": false,
				"extended_tweet": {
					"full_text": "\"The innovative crowdsourcing that the Tagboard, Twitter and TEGNA collaboration enables is surfacing locally relevant conversations in real-time and enabling voters to ask questions during debates,”  -- @adamostrow, @TEGNA\nLearn More: https:\/\/t.co\/ivAFtanfje",
					"display_text_range": [
						0,
						259
					],
					"entities": {
						"hashtags": [],
						"urls": [
							{
								"url": "https:\/\/t.co\/ivAFtanfje",
								"expanded_url": "https:\/\/blog.tagboard.com\/twitter-and-tagboard-collaborate-to-bring-best-election-content-to-news-outlets-with-tagboard-e85fc864bcf4",
								"display_url": "blog.tagboard.com\/twitter-and-ta…",
								"unwound": {
									"url": "https:\/\/blog.tagboard.com\/twitter-and-tagboard-collaborate-to-bring-best-election-content-to-news-outlets-with-tagboard-e85fc864bcf4",
									"status": 200,
									"title": "Twitter and Tagboard Collaborate to Bring Best Election Content to News Outlets With Tagboard…",
									"description": "By Tyler Singletary, Head of Product, Tagboard"
								},
								"indices": [
									236,
									259
								]
							}
						],
						"user_mentions": [
							{
								"screen_name": "adamostrow",
								"name": "Adam Ostrow",
								"id": 5695942,
								"id_str": "5695942",
								"indices": [
									204,
									215
								]
							},
							{
								"screen_name": "TEGNA",
								"name": "TEGNA",
								"id": 34123003,
								"id_str": "34123003",
								"indices": [
									217,
									223
								]
							}
						],
						"symbols": []
					}
				},
				"quote_count": 0,
				"reply_count": 1,
				"retweet_count": 6,
				"favorite_count": 19,
				"entities": {
					"hashtags": [],
					"urls": [
						{
							"url": "https:\/\/t.co\/w46U5TRTzQ",
							"expanded_url": "https:\/\/twitter.com\/i\/web\/status\/1057384253116289025",
							"display_url": "twitter.com\/i\/web\/status\/1…",
							"indices": [
								117,
								140
							]
						}
					],
					"user_mentions": [],
					"symbols": []
				},
				"favorited": false,
				"retweeted": false,
				"possibly_sensitive": false,
				"filter_level": "low",
				"lang": "en"
			},
			"is_quote_status": false,
			"quote_count": 0,
			"reply_count": 0,
			"retweet_count": 0,
			"favorite_count": 0,
			"entities": {
				"hashtags": [],
				"urls": [],
				"user_mentions": [
					{
						"screen_name": "harmophone",
						"name": "Tyler Singletary",
						"id": 175187944,
						"id_str": "175187944",
						"indices": [
							3,
							14
						]
					}
				],
				"symbols": []
			},
			"favorited": false,
			"retweeted": false,
			"filter_level": "low",
			"lang": "en",
			"matching_rules": [
				{
					"tag": null
				}
			]
		}
	],
	"requestParameters": {
		"maxResults": 500,
		"fromDate": "201811010000",
		"toDate": "201811060000"
	}
}`

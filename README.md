# Tweet Cleanup

A utility to quickly delete old tweets in a one-off run, or periodically as a background daemon.


## Requirements

In order to utilize this, you need to create an application in the Twitter developer portal, and reference the keys and tokens for your app.

## Usage

```
Usage of twitter-cleanup:

Environment variables for Twitter authentication:
CONSUMER_SECRET
CONSUMER_KEY
ACCESS_TOKEN
ACCESS_TOKEN_SECRET

Command line flags:
  -from string
    	Starting date to look back to, in Twitter format (YYYYMMDD) (default "20100101")
  -handle string
    	Twitter username to search. (default "smugcloud")
  -period int
    	The frequency with which to check Twitter (in months) (default 1)
  -start int
    	The number of previous months to preserve on Twitter (anything before will be deleted). (default 12)
```

Today, running this will delete the tweets in the specified timeframe, and wait to run again in 1 month (unless you modify this default value).  

## To-Do

* Add a dry-run to see details about what would be deleted
* Add useful endpoints to query the history of the app, if it has been running as a daemon.
* The endpoints assume a cache, which also needs to be added
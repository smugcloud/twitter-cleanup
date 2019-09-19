# Tweet Cleanup

A Lambda (or any serverless really) based application to delete a timeframe of tweets.  You can utilize things like Cloudwatch Events to run this periodically for a low cost solution.

## Requirements

In order to utilize this, you need to create an application in the Twitter developer portal, and reference the keys and tokens for your app.

## Usage

```
Usage of twitter-cleanup:

Environment variables for Twitter authentication and timefram
CONSUMER_SECRET
CONSUMER_KEY
ACCESS_TOKEN
ACCESS_TOKEN_SECRET
FROM
HANDLE
START
```

## To-Do

* Add a dry-run to see details about what would be deleted
* Add useful endpoints to query the history of the app, if it has been running as a daemon.
* The endpoints assume a cache, which also needs to be added
package client

import (
	"fmt"
)

// Sendtweet sends a tweet to Twitter account using the Twitter API. The response message is a string that
// should includ any user tags or hash tags that were used in the tweet. Twitter will return a message with out it being posted if it is
// a duplicate. https://developer.twitter.com/en/docs/twitter-api/v1/tweets/post-and-engage/api-reference/post-statuses-update
func (c *Client) SendTweet(message string) error {
	_, resp, err := c.TweetBot.Statuses.Update(message, nil)
	if err != nil {
		return fmt.Errorf("failed to send tweet: %w", err)
	}
	if resp.StatusCode > 300 {
		return fmt.Errorf("status code: %d\n %v", resp.StatusCode, resp.Body)
	}
	return nil
}

package client

import (
	"fmt"
)

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

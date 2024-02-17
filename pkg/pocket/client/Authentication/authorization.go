package client

import "fmt"

const authUrl = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"

func (c *Client) AuthorizationUrl(redirectUri string, reqToken string) (string, error) {
	return fmt.Sprintf(authUrl, reqToken, redirectUri), nil
}

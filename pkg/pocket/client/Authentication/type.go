package client

import (
	"net/http"
)

const (
	Host             = "getpocket.com"
	PathRequestToken = "/v3/oauth/request"
	PathAccessToken  = "/v3/oauth/authorize"
	PathMethodAdd    = "/v3/add"
)

type ParametersRequestToken struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectUri string `json:"redirect_uri"`
}

type ParametersAccessToken struct {
	ConsumerKey  string `json:"consumer_key"`
	RequestToken string `json:"code"`
}

type RequestToken struct {
	RequestToken string `json:"code"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserName    string `json:"username"`
}

type ParametersAdd struct {
	Url         string `json:"url"`
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}

// client
type Client struct {
	consumerKey string
	client      http.Client
}

func New(consumerKey string) *Client {
	return &Client{
		consumerKey: consumerKey,
		client:      http.Client{},
	}
}

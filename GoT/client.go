package GoT

import (
	"go-twitter/Utils/Proxy"
	"net/http"
)

type client struct {
	consumerKey    string
	consumerSecret string
	bearerToken    string
	client         *http.Client
}

// Client is an API client for Twitter v2 API.
type Client struct {
	*client
}

func New(bearerToken string) *Client {
	c := &client{
		consumerKey:    "",
		consumerSecret: "",
		bearerToken:    bearerToken,
		client:         Proxy.Proxy(),
	}
	return &Client{
		client: c,
	}
}

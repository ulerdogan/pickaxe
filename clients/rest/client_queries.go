package rest_client

import "net/http"

type Client interface {
	Get(url string, headers http.Header) (*http.Response, error)
	Post(url string, body interface{}, headers http.Header) (*http.Response, error)
	NewPriceAPI() PriceAPI
}

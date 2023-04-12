package rest_client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type restClient struct{
	*http.Client
}

func NewRestClient() Client {
	return &restClient{&http.Client{}}
}

func (c *restClient) Get(url string, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = headers

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *restClient) Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	request.Header = headers

	client := http.Client{}
	return client.Do(request)

}

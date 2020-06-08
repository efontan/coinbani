package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Http interface {
	Get(req *GetRequestBuilder) (interface{}, error)
}

type restClient struct {
	client *http.Client
}

func NewRestClient() *restClient {
	c := &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
			MaxIdleConns:        5,
			MaxConnsPerHost:     10,
		},
		Timeout: 10 * time.Second,
	}
	return &restClient{
		client: c,
	}
}

type GetRequestBuilder struct {
	Url           string
	ParseResponse func(response *http.Response) (interface{}, error)
}

func (c *restClient) Get(req *GetRequestBuilder) (interface{}, error) {
	// fetch from service
	r, err := http.NewRequest(http.MethodGet, req.Url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating HTTP Get request")
	}
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36")

	res, err := c.client.Do(r)
	if err != nil || res == nil {
		return nil, errors.Wrap(err, "fetching response from service")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response status code is %d", res.StatusCode))
	}

	return req.ParseResponse(res)
}

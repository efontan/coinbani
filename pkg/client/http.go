package client

import (
	"fmt"
	"net/http"
	"time"

	"coinbani/pkg/cache"

	"github.com/pkg/errors"
)

type Http interface {
	Get(req *GetRequestBuilder) (interface{}, error)
}

type restClient struct {
	client *http.Client
	cache  cache.Cache
}

func NewRestClientWithCache() *restClient {
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
		cache:  cache.New(),
	}
}

type GetRequestBuilder struct {
	Url             string
	CacheKey        string
	CacheExpiration time.Duration
	ParseResponse   func(response *http.Response) (interface{}, error)
}

func (c *restClient) Get(req *GetRequestBuilder) (interface{}, error) {
	v, found := c.cache.Get(req.CacheKey)
	cachedResponse, ok := v.(*http.Response)
	if !ok {
		return nil, errors.New("casting cached response")
	}

	if found {
		// fetch from cache
		return req.ParseResponse(cachedResponse)
	}

	// fetch from service
	r, err := http.NewRequest(http.MethodGet, req.Url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating HTTP Get request")
	}
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36")

	res, err := c.client.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "fetching response from service")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response status code is %d", res.StatusCode))
	}

	c.cache.Set(req.CacheKey, res, req.CacheExpiration)
	return req.ParseResponse(res)
}

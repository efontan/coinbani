package client

import "net/http"

type Http interface {
	Get(url string) (resp *http.Response, err error)
}

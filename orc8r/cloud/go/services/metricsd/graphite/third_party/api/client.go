// Package graphite provides client to graphite-web-api (http://graphite-api.readthedocs.io/en/latest/api.html).
package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

var WrongUrlError = errors.New("Wrong url")

// RequestError is special error, which not only implements an error interface,
// but also provides access to `Query` field, containing original query which
// cause an error.
type RequestError struct {
	Query string
	Type  string
}

func (e RequestError) Error() string {
	return e.Type
}

// Client is client to `graphite web api` (http://graphite-api.readthedocs.io/en/latest/api.html).
// You can either instantiate it manually by providing `url` and `client` or use a `NewFromString` shortcut.
type Client struct {
	Url    url.URL
	Client *http.Client
}

// NewFromString is a convenient function for constructing `graphite.Client` from url string.
// For example 'https://my-graphite.tld'. NewFromString could return either `graphite.Client`
// instance or `WrongUrlError` error.
func NewFromString(urlString string) (*Client, error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return nil, WrongUrlError
	}
	return &Client{*url, &http.Client{}}, nil
}

func (c *Client) makeRequest(q qsGenerator) ([]byte, error) {
	empty := []byte{}
	response, err := c.Client.Get(c.queryAsString(q))
	if err != nil {
		return empty, c.createError(q, "Request error")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return empty, c.createError(q, "Wrong status code")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return empty, c.createError(q, "Can't read response body")
	}
	return body, nil
}

func (c *Client) createError(r qsGenerator, t string) error {
	return RequestError{
		Type:  t,
		Query: c.queryAsString(r),
	}
}

func (c *Client) queryAsString(r qsGenerator) string {
	return c.Url.String() + r.toQueryString()
}

type qsGenerator interface {
	toQueryString() string
}

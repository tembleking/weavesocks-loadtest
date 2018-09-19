package client

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Client interface {
	Login(user, password string) (err error)
	Get(path string) (err error)
	Post(path string, values url.Values) (err error)
	Delete(path string) (err error)
}

type client struct {
	host   string
	client *http.Client
}

func New(host string) (c *client, err error) {
	if host == "" {
		err = errors.New("hostname is empty")
		return
	}

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	if err != nil {
		err = errors.Wrap(err, "error creating the cookie jar")
		return
	}
	c = &client{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar,
		},
		host: host,
	}

	return
}

func (c *client) Login(user, password string) (err error) {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/login", c.host), nil)
	request.SetBasicAuth(user, password)

	response, err := c.client.Do(request)
	if err != nil {
		err = errors.Wrapf(err, "error logging in to %s/login", c.host)
		return
	}
	defer response.Body.Close()

	return
}

func (c *client) Get(path string) (err error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, path), nil)
	if err != nil {
		err = errors.Wrapf(err, "could not create GET request to endpoint %s/%s", c.host, path)
		return
	}
	request.Close = true

	response, err := c.client.Do(request)
	if err != nil {
		err = errors.Wrapf(err, "error with GET request to %s/%s", c.host, path)
		return
	}
	defer response.Body.Close()
	return
}

func (c *client) Post(path string, values url.Values) (err error) {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.host, path), nil)
	if err != nil {
		err = errors.Wrapf(err, "could not create POST request to endpoint %s/%s", c.host, path)
		return
	}
	request.Close = true
	request.PostForm = values

	response, err := c.client.Do(request)
	if err != nil {
		err = errors.Wrapf(err, "error with POST request to %s/%s", c.host, path)
		return
	}
	defer response.Body.Close()
	return
}

func (c *client) Delete(path string) (err error) {

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", c.host, path), nil)
	if err != nil {
		err = errors.Wrapf(err, "error creating DELETE request to endpoint %s/%s", c.host, path)
		return
	}
	request.Close = true

	response, err := c.client.Do(request)
	if err != nil {
		err = errors.Wrapf(err, "error requesting endpoint %s/%s", c.host, path)
		return
	}
	defer response.Body.Close()

	return
}

package client

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

type Client interface {
	Login(user, password string) (err error)
	Get(path string) (err error)
	Post(path string, values url.Values) (err error)
	Delete(path string) (err error)
}

type client struct {
	host   *url.URL
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

	hostUrl, err := url.Parse(host)

	if err != nil {
		err = errors.Wrap(err, "invalid url")
		return
	}

	c = &client{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar,
		},
		host: hostUrl,
	}

	return
}

func (c *client) Login(user, password string) (err error) {
	loginUrl, _ := c.host.Parse("/login")
	request, _ := http.NewRequest(http.MethodGet, loginUrl.String(), nil)
	request.SetBasicAuth(user, password)

	response, err := c.client.Do(request)
	if err != nil {
		err = errors.Wrapf(err, "error logging in to %s/login", c.host)
		return
	}
	defer response.Body.Close()

	return
}

func (c *client) Get(endpoint string) (err error) {
	requestUrl, _ := c.host.Parse(fmt.Sprintf("/%s", endpoint))
	request, err := http.NewRequest(http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		err = errors.Wrapf(err, "could not create GET request to endpoint %s", requestUrl.String())
		return
	}
	request.Close = true

	response, err := c.client.Do(request)
	if err != nil {
		err = errors.Wrapf(err, "error with GET request to %s/%s", c.host, endpoint)
		return
	}
	defer response.Body.Close()
	return
}

func (c *client) Post(endpoint string, values url.Values) (err error) {
	requestUrl, _ := c.host.Parse(fmt.Sprintf("/%s", endpoint))
	request, err := http.NewRequest(http.MethodPost, requestUrl.String(), nil)
	if err != nil {
		err = errors.Wrapf(err, "could not create POST request to endpoint %s", requestUrl.String())
		return
	}
	request.Close = true
	request.PostForm = values

	response, err := c.client.Do(request)
	if err != nil {
		err = errors.Wrapf(err, "error with POST request to %s", requestUrl.String())
		return
	}
	defer response.Body.Close()
	return
}

func (c *client) Delete(endpoint string) (err error) {
	requestUrl, _ := c.host.Parse(fmt.Sprintf("/%s", endpoint))
	request, err := http.NewRequest(http.MethodDelete, requestUrl.String(), nil)
	if err != nil {
		err = errors.Wrapf(err, "error creating DELETE request to endpoint %s", requestUrl.String())
		return
	}
	request.Close = true

	response, err := c.client.Do(request)
	if err != nil {
		err = errors.Wrapf(err, "error requesting endpoint %s", requestUrl.String())
		return
	}
	defer response.Body.Close()

	return
}

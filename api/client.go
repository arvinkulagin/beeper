package api

import (
	"net/http"
	"net/url"
	"bytes"
	"errors"
	"encoding/json"
	"io/ioutil"
)

const (
	Scheme = "http"
	TopicPrefix = "topic"
	PingPrefix = "ping"
)

type Client struct {
	URL url.URL
}

func NewClient(host string) (Client, error) {
	c := Client{
		URL: url.URL{
			Scheme: Scheme,
			Host: host,
		},
	}
	if err := c.Ping(); err != nil {
		return c, err
	}
	return c, nil
}

func (c Client) Add(topic string) error {
	out := bytes.NewReader([]byte(topic))
	c.URL.Path = c.URL.Path + TopicPrefix
	request, err := http.NewRequest("POST", c.URL.String(), out)
	if err != nil {
		return errors.New("Bad request")
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New("Сan not connect to " + c.URL.Host)
	}
	if response.StatusCode != 200 {
		return errors.New("Topic already exists")
	}
	return nil
}

func (c Client) Del(topic string) error {
	c.URL.Path = c.URL.Path + TopicPrefix + "/" + topic
	req, err := http.NewRequest("DELETE", c.URL.String(), nil)
	if err != nil {
		return errors.New("Can not make request")
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Сan not connect to " + c.URL.Host)
	}
	if r.StatusCode != 200 {
		return errors.New("Topic does not exist")
	}
	return nil
}

func (c Client) Pub(topic string, data string) error {
	out := bytes.NewReader([]byte(data))
	c.URL.Path = c.URL.Path + TopicPrefix + "/" + topic
	request, err := http.NewRequest("POST", c.URL.String(), out)
	if err != nil {
		return errors.New("Bad request")
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New("Сan not connect to " + c.URL.Host)
	}
	if response.StatusCode != 200 {
		return errors.New("Topic does not exist")
	}
	return nil
}

func (c Client) List() ([]string, error) {
	list := []string{}
	c.URL.Path = c.URL.Path + TopicPrefix
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		return list, errors.New("Can not make request")
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return list, errors.New("Сan not connect to " + c.URL.Host)
	}
	if r.StatusCode == 500 {
		return list, errors.New("Server error")
	}
	if r.StatusCode != 200 {
		return list, errors.New("Topic does not exist")
	}
	body, err := ioutil.ReadAll(r.Body) // Ограничить количество байт из ReadAll с помощью io.LimitReader
	defer r.Body.Close()
	if err != nil {
		return list, errors.New("Can not read response body")
	}
	json.Unmarshal(body, &list)
	return list, nil
}

func (c Client) Ping() error {
	c.URL.Path = c.URL.Path + PingPrefix
	request, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		return errors.New("Bad request")
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New("Сan not connect to " + c.URL.Host)
	}
	if response.StatusCode != 200 {
		return errors.New("Service is unavailable")
	}
	return nil
}
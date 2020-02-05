package agilecrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var ErrInvalidCreds = fmt.Errorf("username or password is invalid")
var ErrUnauthorized = fmt.Errorf("unauthorized")
var ErrWrongFormat = fmt.Errorf("input in wrong format")
var ErrContactLimit = fmt.Errorf("limit of contacts exceeded")
var ErrNoContacts = fmt.Errorf("no contacts in account")
var ErrNoSuchContact = fmt.Errorf("no contact with that ID found")

const apiURLf = "https://%v.agilecrm.com/dev/"

// TODO: hoist all route bases to here to make them consts

type Client struct {
	url string
	ht  http.Client
}

// route ...
func (c *Client) route(r string) string {
	return fmt.Sprintf("%v%v", c.url, r)
}

type Config struct {
	Domain        string
	User          string
	Password      string
	DefaultClient http.RoundTripper
}

// New ...
func New(conf Config) (*Client, error) {
	if conf.Domain == "" {
		return nil, fmt.Errorf("domain is required in config")
	}

	url := fmt.Sprintf(apiURLf, conf.Domain)

	cl, err := getClient(conf)
	if err != nil {
		return nil, err
	}

	return &Client{url: url, ht: cl}, nil
}

func getClient(conf Config) (http.Client, error) {
	user := strings.TrimSpace(conf.User)
	pass := strings.TrimSpace(conf.Password)

	if user == "" || pass == "" {
		return http.Client{}, ErrInvalidCreds
	}

	return http.Client{
		Transport: rt{
			user: user,
			pass: pass,
			orig: conf.DefaultClient,
		},
		Timeout: time.Second * 10,
	}, nil
}

// req ...
func (c *Client) req(method, route string, body io.Reader, params map[string]string) (*http.Request, error) {
	r := c.route(route)
	req, err := http.NewRequest(method, r, body)
	if err != nil {
		return nil, err
	}

	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("Accept", "application/json")

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	return req, nil
}

// postForm ...
func (c *Client) postForm(method, route string, body io.Reader, params map[string]string) (*http.Request, error) {
	req, err := c.req(method, route, body, params)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

// _processResults ...
func (c *Client) processResults(req *http.Request, out interface{}) (int, error) {
	res, err := c.ht.Do(req)
	if err != nil {
		return -1, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return -1, err
	}

	// fmt.Printf("status: %v\n", res.StatusCode)
	// fmt.Printf("body: \n\n%v\n\n", string(resBody))

	if len(resBody) <= 0 || res.StatusCode == http.StatusNoContent {
		return res.StatusCode, nil
	}

	err = json.Unmarshal(resBody, out)
	if err != nil {
		return -1, err
	}

	return res.StatusCode, nil
}

// _getContact ...
func (c *Client) get(method, route string, body io.Reader, params map[string]string, out interface{}) (int, error) {
	req, err := c.req(method, route, body, params)
	if err != nil {
		return -1, err
	}

	return c.processResults(req, out)
}

// findByID ...
func (c *Client) findByID(route string, out interface{}) error {
	st, err := c.get("GET", route, nil, nil, &out)
	if st == http.StatusOK && err == nil {
		return nil
	}

	switch st {
	case http.StatusNoContent:
		return ErrNoSuchContact
	case http.StatusUnauthorized:
		return ErrUnauthorized
	}

	return err
}

// _sendContact ...
func (c *Client) send(method, route string, params map[string]string, in interface{}, out interface{}) (int, error) {
	bits, err := json.Marshal(in)
	if err != nil {
		return -1, err
	}
	buf := bytes.NewBuffer(bits)
	return c.get(method, route, buf, params, out)
}

// delete ...
func (c *Client) delete(route string) error {
	req, err := c.req("DELETE", route, nil, nil)
	if err != nil {
		return err
	}

	res, err := c.ht.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 204 {
		return fmt.Errorf(res.Status)
	}
	return nil
}

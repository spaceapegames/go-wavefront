// Package wavefront provides a library for interacting with the Wavefront API,
// along with a writer for sending metrics
package wavefront

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Config struct {
	Address       string
	Token         string
	SkipTLSVerify bool // Dev/Test only
}

type Client struct {
	Config     *Config
	BaseURL    *url.URL
	httpClient *http.Client
	debug      bool

	Alerts Alerting
	Query  Querying
	Events Event
}

type QueryParams map[string]string

// NewClient returns a new Wavefront client
func NewClient(config *Config) (*Client, error) {
	baseURL, err := url.Parse("https://" + config.Address)
	if err != nil {
		return nil, err
	}

	c := &Client{Config: config,
		BaseURL:    baseURL,
		httpClient: http.DefaultClient,
		debug:      false,
	}

	//For testing ONLY
	if config.SkipTLSVerify == true {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.httpClient.Transport = tr
	}

	c.Alerts = Alerting{client: c}
	c.Query = Querying{client: c}
	c.Events = Event{client: c}

	return c, nil
}

// NewRequest creates a request object to query Wavefront
// A relative URL should be passed along with the method, that will be resolved against the client's BaseURL
// Parameters can be passed, a QueryParams map of k: v
func (client Client) NewRequest(method, path string, params *QueryParams) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	url := client.BaseURL.ResolveReference(rel)

	if params != nil {
		q := url.Query()
		for k, v := range *params {
			q.Set(k, v)
		}
		url.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-AUTH-TOKEN", client.Config.Token)
	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
}

// Do executes a request against the Wavefront API
// and decodes the JSON into parseInto.
func (client Client) Do(req *http.Request, parseInto interface{}) (io.Reader, error) {

	if client.debug == true {
		d, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%s\n", d)
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Server returned %s\n", resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	// bytes.Reader implements Seek, which we need to use to 'rewind' the Body
	r := bytes.NewReader(body)

	resp.Body.Close()

	if err := json.NewDecoder(r).Decode(&parseInto); err != nil {
		return nil, err
	}

	// 'rewind' the raw response, to make it useful
	r.Seek(0, 0)
	return r, nil
}

// Debug enables dumping http request objects to stdout
func (client *Client) Debug(enable bool) {
	client.debug = enable
}

package api

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/ghodss/yaml"

	"github.com/digitalrebar/provision/models"
)

const APIPATH = "/api/v3"

type Decoder interface {
	Decode(interface{}) error
}

type Encoder interface {
	Encode(interface{}) error
}

func DecodeYaml(buf []byte, ref interface{}) error {
	return yaml.Unmarshal(buf, ref)
}

func Unmarshal(resp *http.Response, ref interface{}) error {
	var dec Decoder
	if resp != nil {
		defer resp.Body.Close()
	}
	ct := resp.Header.Get("Content-Type")
	mt, _, _ := mime.ParseMediaType(ct)
	switch mt {
	case "application/json":
		dec = json.NewDecoder(resp.Body)
	case "application/octet-stream":
		if wr, ok := ref.(io.Writer); !ok {
			return fmt.Errorf("Response is an octet stream, expected ref to be an io.Writer")
		} else {
			_, err := io.Copy(wr, resp.Body)
			return err
		}
	default:
		return fmt.Errorf("Cannot handle content-type %s", ct)
	}
	if dec == nil {
		return fmt.Errorf("No decoder for content-type %s", ct)
	}
	if resp.StatusCode >= 400 {
		res := &models.Error{}
		if err := dec.Decode(res); err != nil {
			return err
		}
		return res
	}
	return dec.Decode(ref)
}

type Client struct {
	*http.Client
	endpoint, username, password string
	token                        *models.UserToken
	closer                       chan struct{}
	closed                       bool
}

func (c *Client) Close() {
	c.closer <- struct{}{}
	c.closed = true
}

func (c *Client) Token() string {
	return c.token.Token
}

func (c *Client) Info() *models.Info {
	return &c.token.Info
}

func (c *Client) UrlFor(args ...string) *url.URL {
	res, err := url.ParseRequestURI(c.endpoint + path.Join(APIPATH, path.Join(args...)))
	if err != nil {
		log.Panicf("Unable to form URL for %v\n    %v", args, err)
	}
	return res
}

func (c *Client) Authorize(req *http.Request) error {
	req.Header.Add("Authorization", "Bearer "+c.Token())
	return nil
}

func (c *Client) Request(method string, uri *url.URL, body io.Reader) (*http.Request, error) {
	if c.closed {
		return nil, fmt.Errorf("Connection closed")
	}
	req, err := http.NewRequest(method, uri.String(), body)
	if err == nil {
		err = c.Authorize(req)
	}
	return req, err
}

func (c *Client) RequestJSON(method string, uri *url.URL, body io.Reader) (*http.Request, error) {
	req, err := c.Request(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) DoJSON(method string, uri *url.URL, body io.Reader, val interface{}) error {
	req, err := c.RequestJSON(method, uri, body)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if val != nil {
		return Unmarshal(resp, val)
	}
	return nil
}

func (c *Client) ListBlobs(at string, params map[string]string) ([]string, error) {
	reqURI := c.UrlFor(path.Join("/", at))
	vals := url.Values{}
	for k, v := range params {
		vals.Add(k, v)
	}
	reqURI.RawQuery = vals.Encode()
	res := []string{}
	return res, c.DoJSON("GET", reqURI, nil, res)
}

func (c *Client) GetBlob(at ...string) (io.ReadCloser, error) {
	reqURI := c.UrlFor(path.Join("/", path.Join(at...)))
	req, err := c.Request("GET", reqURI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Add("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) PostBlob(blob io.Reader, at ...string) error {
	reqURI := c.UrlFor(path.Join("/", path.Join(at...)))
	req, err := c.Request("POST", reqURI, blob)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) DeleteBlob(at ...string) error {
	reqURI := c.UrlFor(path.Join("/", path.Join(at...)))
	req, err := c.RequestJSON("DELETE", reqURI, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) AllIndexes() (map[string]map[string]models.Index, error) {
	reqURI := c.UrlFor("indexes")
	res := map[string]map[string]models.Index{}
	return res, c.DoJSON("GET", reqURI, nil, &res)
}

func (c *Client) Indexes(prefix string) (map[string]models.Index, error) {
	reqURI := c.UrlFor("indexes", prefix)
	res := map[string]models.Index{}
	return res, c.DoJSON("GET", reqURI, nil, &res)
}

func (c *Client) OneIndex(prefix, param string) (models.Index, error) {
	reqURI := c.UrlFor("indexes", prefix, param)
	res := models.Index{}
	return res, c.DoJSON("GET", reqURI, nil, &res)
}

func (c *Client) ListModels(ref models.Models, params map[string]string) error {
	reqURI := c.UrlFor(ref.Elem().Prefix())
	vals := url.Values{}
	for k, v := range params {
		vals.Add(k, v)
	}
	reqURI.RawQuery = vals.Encode()
	return c.DoJSON("GET", reqURI, nil, ref)
}

func (c *Client) GetModel(prefix, key string) (models.Model, error) {
	res, err := models.New(prefix)
	if err != nil {
		return nil, err
	}
	reqURI := c.UrlFor(prefix, key)
	return res, c.DoJSON("GET", reqURI, nil, res)
}

func (c *Client) FillModel(ref models.Model, key string) error {
	reqURI := c.UrlFor(ref.Prefix(), key)
	return c.DoJSON("GET", reqURI, nil, ref)
}

func (c *Client) CreateModel(ref models.Model) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(ref); err != nil {
		return err
	}
	reqURI := c.UrlFor(ref.Prefix())
	return c.DoJSON("POST", reqURI, buf, ref)
}

func (c *Client) DeleteModel(prefix, key string) (models.Model, error) {
	res, err := models.New(prefix)
	if err != nil {
		return nil, err
	}
	reqURI := c.UrlFor(prefix, key)
	return res, c.DoJSON("POST", reqURI, nil, res)
}

func (c *Client) reauth(tok *models.UserToken) error {
	reqURI := c.UrlFor("users", c.username, "token")
	v := url.Values{}
	v.Set("ttl", "600")
	reqURI.RawQuery = v.Encode()
	return c.DoJSON("GET", reqURI, nil, tok)
}

func (c *Client) PatchModel(old, new models.Model) error {
	if old.Prefix() != new.Prefix() {
		log.Panicf("Cannot patch %s into a %s", old.Prefix(), new.Prefix())
	}
	oldBuf, err := json.Marshal(old)
	if err != nil {
		return err
	}
	newBuf, err := json.Marshal(new)
	if err != nil {
		return err
	}
	patch, err := jsonpatch2.Generate(oldBuf, newBuf, true)
	if err != nil {
		return err
	}
	buf, err := json.Marshal(patch)
	reqURI := c.UrlFor(old.Prefix(), old.Key())
	return c.DoJSON("PATCH", reqURI, bytes.NewBuffer(buf), new)
}

func (c *Client) PutModel(obj models.Model) error {
	reqURI := c.UrlFor(obj.Prefix(), obj.Key())
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return err
	}
	return c.DoJSON("PUT", reqURI, buf, obj)
}

func TokenSession(endpoint, token string) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &Client{
		endpoint: endpoint,
		Client:   &http.Client{Transport: tr},
		closer:   make(chan struct{}, 0),
		token:    &models.UserToken{Token: token},
	}
	return c, nil
}

func UserSession(endpoint, username, password string) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &Client{
		endpoint: endpoint,
		username: username,
		password: password,
		Client:   &http.Client{Transport: tr},
		closer:   make(chan struct{}, 0),
	}
	req, err := c.RequestJSON("GET", c.UrlFor("users", c.username, "token"), nil)
	if err != nil {
		return nil, err
	}
	basicAuth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Set("Authorization", "Basic "+basicAuth)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	token := &models.UserToken{}
	if err := Unmarshal(resp, token); err != nil {
		return nil, err
	}
	go func() {
		ticker := time.NewTicker(300 * time.Second)
		for {
			select {
			case <-c.closer:
				ticker.Stop()
				return
			case <-ticker.C:
				token := &models.UserToken{}
				if err := c.reauth(token); err != nil {
					log.Fatalf("Error reauthing token, aborting: %v", err)
				}
				c.token = token
			}
		}
	}()
	c.token = token
	return c, nil
}

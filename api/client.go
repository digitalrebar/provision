// Package api implements a client API for working with
// digitalrebar/provision.
package api

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/VictorLowther/jsonpatch2"

	"github.com/digitalrebar/provision/models"
)

// APIPATH is the base path for all API endpoints that digitalrebar
// provision provides.
const APIPATH = "/api/v3"

// Client wraps *http.Client to include our authentication routines
// and routines for handling some of the biolerplate CRUD operations
// against digitalrebar provision.
type Client struct {
	*http.Client
	endpoint, username, password string
	token                        *models.UserToken
	closer                       chan struct{}
	closed                       bool
}

// Close should be called whenever you no longer want to use this
// client connection.  It will stop any token refresh routines running
// in the background, and force any API calls made to this client that
// would communicate with the server to return an error
func (c *Client) Close() {
	c.closer <- struct{}{}
	close(c.closer)
	c.closed = true
}

// Token returns the current authentication token associated with the
// Client.
func (c *Client) Token() string {
	if c.token == nil {
		return ""
	}
	return c.token.Token
}

// Info returns some basic system information that was retrieved as
// part of the initial authentication.
func (c *Client) Info() (*models.Info, error) {
	res := &models.Info{}
	return res, c.DoJSON("GET", c.UrlFor("info"), nil, res)
}

// UrlFor is a helper function used to build URLs for the other client
// helper functions.
func (c *Client) UrlFor(args ...string) *url.URL {
	res, err := url.ParseRequestURI(c.endpoint + path.Join(APIPATH, path.Join(args...)))
	if err != nil {
		log.Panicf("Unable to form URL for %v\n    %v", args, err)
	}
	return res
}

func (c *Client) WithParams(uri *url.URL, params map[string]string) {
	if params != nil && len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		uri.RawQuery = values.Encode()
	}
}

// Authorize sets the Authorization header in the Request with the
// current bearer token.  The rest of the helper methods call this, so
// you don't have to unless you are building your own http.Requests.
func (c *Client) Authorize(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+c.Token())
	return nil
}

// Request builds a preauthorized http.Request.
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

// RequestJSON builds an http.Request that has the Accept and
// Content-Type headers set to application/json
func (c *Client) RequestJSON(method string, uri *url.URL, body io.Reader) (*http.Request, error) {
	req, err := c.Request(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) doJSON(method string, uri *url.URL, body io.Reader, val interface{}) error {
	req, err := c.RequestJSON(method, uri, body)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	return Unmarshal(resp, val)
}

// DoJSON does a complete round-trip, unmarshalling the body returned
// from the server into val.  It calls RequestJSON to build the
// request.
func (c *Client) DoJSON(method string, uri *url.URL, body interface{}, val interface{}) error {
	switch obj := body.(type) {
	case nil:
		return c.doJSON(method, uri, nil, val)
	case io.Reader:
		return c.doJSON(method, uri, obj, val)
	case []byte:
		return c.doJSON(method, uri, bytes.NewBuffer(obj), val)
	default:
		buf, err := json.Marshal(&body)
		if err != nil {
			return err
		}
		return c.doJSON(method, uri, bytes.NewBuffer(buf), val)
	}
}

// ListBlobs lists the names of all the binary objects at 'at', using
// the indexing parameters suppied by params.
func (c *Client) ListBlobs(at string, params map[string]string) ([]string, error) {
	reqURI := c.UrlFor(path.Join("/", at))
	c.WithParams(reqURI, params)
	res := []string{}
	return res, c.DoJSON("GET", reqURI, nil, &res)
}

// GetBlob fetches a binary blob from the server.  You are responsible
// for copying the returned io.ReadCloser to a suitable location and
// closing it afterwards if it is not nil, otherwise the client will
// leak open HTTP connections.
func (c *Client) GetBlob(dest io.Writer, at ...string) error {
	reqURI := c.UrlFor(path.Join("/", path.Join(at...)))
	req, err := c.Request("GET", reqURI, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Add("Accept", "application/json")
	resp, err := c.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	return Unmarshal(resp, dest)
}

// PostBlob uploads the binary blob contained in the passed io.Reader
// to the location specified by at on the server.  You are responsible
// for closing the passed io.Reader.
func (c *Client) PostBlob(blob io.Reader, at ...string) (models.BlobInfo, error) {
	reqURI := c.UrlFor(path.Join("/", path.Join(at...)))
	res := models.BlobInfo{}
	req, err := c.Request("POST", reqURI, blob)
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	return res, Unmarshal(resp, &res)
}

// DeleteBlob deletes a blob on the server at the location indicated
// by 'at'
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
	defer resp.Body.Close()
	return Unmarshal(resp, nil)
}

// AllIndexes returns all the static indexes available for all object
// types on the server.
func (c *Client) AllIndexes() (map[string]map[string]models.Index, error) {
	reqURI := c.UrlFor("indexes")
	res := map[string]map[string]models.Index{}
	return res, c.DoJSON("GET", reqURI, nil, &res)
}

// Indexes returns all the static indexes available for a given type
// of object on the server.
func (c *Client) Indexes(prefix string) (map[string]models.Index, error) {
	reqURI := c.UrlFor("indexes", prefix)
	res := map[string]models.Index{}
	return res, c.DoJSON("GET", reqURI, nil, &res)
}

// OneIndex tests to see if there is an index on the object type
// indicated by prefix for a specific parameter.  If the returned
// Index is empty, there is no such Index.
func (c *Client) OneIndex(prefix, param string) (models.Index, error) {
	reqURI := c.UrlFor("indexes", prefix, param)
	res := models.Index{}
	return res, c.DoJSON("GET", reqURI, nil, &res)
}

func (c *Client) ListModel(prefix string, params map[string]string) ([]models.Model, error) {
	ref, err := models.New(prefix)
	if err != nil {
		return nil, err
	}
	reqURI := c.UrlFor(ref.Prefix())
	c.WithParams(reqURI, params)
	res := ref.SliceOf()
	if err := c.DoJSON("GET", reqURI, nil, res); err != nil {
		return nil, err
	}
	return ref.ToModels(res), nil
}

// GetModel returns an object if type prefix with the unique
// identifier key, if such an object exists.  Key can be either the
// unique ket for an object, or any field on an object that has an
// index that enforces uniqueness.
func (c *Client) GetModel(prefix, key string) (models.Model, error) {
	res, err := models.New(prefix)
	if err != nil {
		return nil, err
	}
	reqURI := c.UrlFor(res.Prefix(), key)
	return res, c.DoJSON("GET", reqURI, nil, &res)
}

// ExistsModel tests to see if an object exists on the server
// following the same rules as GetModel
func (c *Client) ExistsModel(prefix, key string) (bool, error) {
	reqURI := c.UrlFor(prefix, key)
	req, err := c.Request("HEAD", reqURI, nil)
	if err != nil {
		return false, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return false, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		res := &models.Error{Code: resp.StatusCode, Type: prefix, Key: key}
		res.Errorf("Unable to determine existence")
		return false, res
	}
}

// FillModel fills the passed-in model with new information retrieved
// from the server.
func (c *Client) FillModel(ref models.Model, key string) error {
	reqURI := c.UrlFor(ref.Prefix(), key)
	err := c.DoJSON("GET", reqURI, nil, &ref)
	if f, ok := ref.(models.Filler); err == nil && ok {
		f.Fill()
	}
	return err
}

// CreateModel takes the passed-in model and creates an instance of it
// on the server.  It will return an error if the passed-in model does
// not validate or if it already exists on the server.
func (c *Client) CreateModel(ref models.Model) error {
	reqURI := c.UrlFor(ref.Prefix())
	if f, ok := ref.(models.Filler); ok {
		f.Fill()
	}
	return c.DoJSON("POST", reqURI, ref, &ref)
}

// DeleteModel deletes the model matching the passed-in prefix and
// key.  It returns the object that was deleted.
func (c *Client) DeleteModel(prefix, key string) (models.Model, error) {
	res, err := models.New(prefix)
	if err != nil {
		return nil, err
	}
	reqURI := c.UrlFor(prefix, key)
	return res, c.DoJSON("DELETE", reqURI, nil, &res)
}

func (c *Client) reauth(tok *models.UserToken) error {
	reqURI := c.UrlFor("users", c.username, "token")
	v := url.Values{}
	v.Set("ttl", "600")
	reqURI.RawQuery = v.Encode()
	return c.DoJSON("GET", reqURI, nil, tok)
}

// PatchModel attempts to update the object matching the passed prefix
// and key on the server side with the passed-in JSON patch (as
// sepcified in https://tools.ietf.org/html/rfc6902).  To ensure that
// conflicting changes are rejected, your patch should contain the
// appropriate test stanzas, which will allow the server to detect and
// reject conflicting changes from different sources.
func (c *Client) PatchModel(prefix, key string, patch jsonpatch2.Patch) (models.Model, error) {
	new, err := models.New(prefix)
	if err != nil {
		return nil, err
	}
	reqURI := c.UrlFor(prefix, key)
	err = c.DoJSON("PATCH", reqURI, patch, &new)
	if err == nil {
		new.Fill()
	}
	return new, err
}

// PutModel replaces the server-side object matching the passed-in
// object with the passed-in object.  Note that PutModel does not
// allow the server to detect and reject conflicting changes from
// multiple sources.
func (c *Client) PutModel(obj models.Model) error {
	reqURI := c.UrlFor(obj.Prefix(), obj.Key())
	return c.DoJSON("PUT", reqURI, obj, &obj)
}

// TokenSession creates a new api.Client that will use the passed-in Token for authentication.
// It should be used whenever the API is not acting on behalf of a user.
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
	go func() {
		<-c.closer
	}()
	return c, nil
}

// UserSession creates a new api.Client that can act on behalf of a
// user.  It will perform a single request using basic authentication
// to get a token that expires 600 seconds from the time the session
// is crated, and every 300 seconds it will refresh that token.
//
// UserSession does not currently attempt to cache tokens to
// persistent storage, although that may change in the future.
func UserSession(endpoint, username, password string) (*Client, error) {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
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

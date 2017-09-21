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
	"mime"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/ghodss/yaml"

	"github.com/digitalrebar/provision/models"
)

// APIPATH is the base path for all API endpoints that digitalrebar
// provision provides.
const APIPATH = "/api/v3"

type Decoder interface {
	Decode(interface{}) error
}

type Encoder interface {
	Encode(interface{}) error
}

// DecodeYaml is a helper function for dealing with user input -- when
// accepting input from the user, we want to treat both YAML and JSON
// as first-class citizens.  The YAML library we use makes that easier
// by using the json struct tags for all marshalling and unmarshalling
// purposes.
//
// Note that the REST API does not use YAML as a wire protocol, so
// this function should never be used to decode data coming from the
// provision service.
func DecodeYaml(buf []byte, ref interface{}) error {
	return yaml.Unmarshal(buf, ref)
}

// Unmarshal is a helper for decoding the body of a response from the server.
// It should be called in one of two ways:
//
// The first is when you expect the response body to contain a blob
// of data that needs to be streamed somewhere.  In that case, ref
// should be an io.Writer, and the Content-Type header will be ignored.
//
// The second is when you expect the response body to contain a
// serialized object to be unmarshalled.  In that case, the response's
// Content-Type will be used as a hint to decide how to unmarshall the
// recieved data into ref.
//
// In either case, if there are any errors in the unmarshalling
// process or the response StatusCode indicates non-success, an error
// will be returned and you should not expect ref to contain vaild
// data.
func Unmarshal(resp *http.Response, ref interface{}) error {
	if resp != nil {
		defer resp.Body.Close()
	}
	if wr, ok := ref.(io.Writer); ok && resp.StatusCode < 300 {
		_, err := io.Copy(wr, resp.Body)
		return err
	}
	var dec Decoder
	ct := resp.Header.Get("Content-Type")
	mt, _, _ := mime.ParseMediaType(ct)
	switch mt {
	case "application/json":
		dec = json.NewDecoder(resp.Body)
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
	return c.token.Token
}

// Info returns some basic system information that was retrieved as
// part of the initial authentication.
func (c *Client) Info() *models.Info {
	return &c.token.Info
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

// DoJSON does a complete round-trip, unmarshalling the body returned
// from the server into val.  It calls RequestJSON to build the
// request.
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

// ListBlobs lists the names of all the binary objects at 'at', using
// the indexing parameters suppied by params.
func (c *Client) ListBlobs(at string, params map[string]string) ([]string, error) {
	reqURI := c.UrlFor(path.Join("/", at))
	if params != nil {
		vals := url.Values{}
		for k, v := range params {
			vals.Add(k, v)
		}
		reqURI.RawQuery = vals.Encode()
	}
	res := []string{}
	return res, c.DoJSON("GET", reqURI, nil, res)
}

// GetBlob fetches a binary blob from the server.  You are responsible
// for copying the returned io.ReadCloser to a suitable location and
// closing it afterwards if it is not nil, otherwise the client will
// leak open HTTP connections.
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

// PostBlob uploads the binary blob contained in the passed io.Reader
// to the location specified by at on the server.  You are responsible
// for closing the passed io.Reader.
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
	resp.Body.Close()
	return nil
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

// ListModels returns all of the objects matching the passed params.
// If no params are passed, all objects of the specified type are
// returned.
func (c *Client) ListModels(ref models.Models, params map[string]string) error {
	reqURI := c.UrlFor(ref.Elem().Prefix())
	if params != nil {
		vals := url.Values{}
		for k, v := range params {
			vals.Add(k, v)
		}
		reqURI.RawQuery = vals.Encode()
	}
	return c.DoJSON("GET", reqURI, nil, ref)
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
	reqURI := c.UrlFor(prefix, key)
	return res, c.DoJSON("GET", reqURI, nil, res)
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
	return c.DoJSON("GET", reqURI, nil, ref)
}

// CreateModel takes the passed-in model and creates an instance of it
// on the server.  It will return an error if the passed-in model does
// not validate or if it already exists on the server.
func (c *Client) CreateModel(ref models.Model) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(ref); err != nil {
		return err
	}
	reqURI := c.UrlFor(ref.Prefix())
	return c.DoJSON("POST", reqURI, buf, ref)
}

// DeleteModel deletes the model matching the passed-in prefix and
// key.  It returns the object that was deleted.
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

// PatchModel attempts to update the object matching the passed prefix
// and key on the server side with the passed-in JSON patch (as
// sepcified in https://tools.ietf.org/html/rfc6902).  To ensure that
// conflicting changes are rejected, your patch should contain the
// appropriate test stanzas, which will allow the server to detect and
// reject conflicting changes from different sources.
func (c *Client) PatchModel(prefix, key string, patch *jsonpatch2.Patch) (models.Model, error) {
	new, err := models.New(prefix)
	if err != nil {
		return nil, err
	}
	buf, err := json.Marshal(patch)
	reqURI := c.UrlFor(prefix, key)
	return new, c.DoJSON("PATCH", reqURI, bytes.NewBuffer(buf), new)
}

// PutModel replaces the server-side object matching the passed-in
// object with the passed-in object.  Note that PutModel does not
// allow the server to detect and reject conflicting changes from
// multiple sources.
func (c *Client) PutModel(obj models.Model) error {
	reqURI := c.UrlFor(obj.Prefix(), obj.Key())
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return err
	}
	return c.DoJSON("PUT", reqURI, buf, obj)
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

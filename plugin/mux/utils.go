package mux

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	rt "runtime/debug"

	"github.com/digitalrebar/logger"
	"github.com/digitalrebar/provision/models"
)

// JsonResponse returns a JSON object on the http writer.
// Setting the code and encoding the provided object.
func JsonResponse(w http.ResponseWriter, code int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	return enc.Encode(obj)
}

// TestContentType tests to ensure that requested content type is
// in the header.
func TestContentType(r *http.Request, ct string) (string, bool) {
	rct := r.Header.Get("Content-Type")
	return rct, strings.Contains(
		strings.ToUpper(rct),
		strings.ToUpper(ct))
}

// AssureContentType returns true if the requested content-type is
// present in the request header and returns a JSON encoded models.Error
// on the http stream and false value if it is not.
func AssureContentType(w http.ResponseWriter, r *http.Request, ct string) bool {
	rct, ok := TestContentType(r, ct)
	if ok {
		return true
	}
	err := &models.Error{Type: r.Method, Code: http.StatusBadRequest}
	err.Errorf("Invalid content type: %s", rct)
	JsonResponse(w, err.Code, err)
	return false
}

// AssureDecode returns true if can successfully decode the incoming JSON
// object into the provided interface.  It will return false and generate
// a JSON encoded error message upon failure.
func AssureDecode(w http.ResponseWriter, r *http.Request, val interface{}) bool {
	if !AssureContentType(w, r, "application/json") {
		return false
	}
	if r.ContentLength == 0 || r.Body == nil {
		val = nil
		return true
	}
	dec := json.NewDecoder(r.Body)
	marshalErr := dec.Decode(&val)
	if marshalErr == nil {
		return true
	}
	err := &models.Error{Type: r.Method, Code: http.StatusBadRequest}
	err.AddError(marshalErr)
	JsonResponse(w, err.Code, err)
	return false
}

// Get provides a path for the Plugin to get data from
// DRP over the Plugin Server RestFUL API.
func Get(client *http.Client, path string) ([]byte, error) {
	resp, err := client.Get(fmt.Sprintf("http://unix/api-server-plugin/v3%s", path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}
	if resp.StatusCode >= 400 {
		berr := &models.Error{}
		err := json.Unmarshal(b, berr)
		if err != nil {
			return nil, e
		}
		return nil, berr
	}
	return b, nil
}

// Post provides a path for the Plugin to send messages to
// DRP over the Plugin Server RestFUL API.
func Post(client *http.Client, path string, indata interface{}) ([]byte, error) {
	data, err := json.Marshal(indata)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(
		fmt.Sprintf("http://unix/api-server-plugin/v3%s", path),
		"application/json",
		strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}
	if resp.StatusCode >= 400 {
		berr := &models.Error{}
		err := json.Unmarshal(b, berr)
		if err != nil {
			return nil, e
		}
		return nil, berr
	}
	return b, nil
}

// Delete provides a path for the Plugin to delete data from
// DRP over the Plugin Server RestFUL API.
func Delete(client *http.Client, path string) error {
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("http://unix/api-server-plugin/v3%s", path),
		nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return e
	}
	if resp.StatusCode >= 400 {
		berr := &models.Error{}
		err := json.Unmarshal(b, berr)
		if err != nil {
			return e
		}
		return berr
	}
	return nil
}

// ResponseWriter wraps a normal http.ResponseWriter with
// a logger and a status value.
type ResponseWriter struct {
	http.ResponseWriter
	logger.Logger
	status int
}

// Status returns the current status ofthe ResponseWriter.
func (rw *ResponseWriter) Status() int {
	return rw.status
}

// Write sends the specified buf after setting StatusOK
// if no other status has been written.
func (rw *ResponseWriter) Write(buf []byte) (int, error) {
	if rw.status == 0 {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(buf)
}

// WriteHeader takes the status code and writes the header
// out the http stream after recording the status.
func (rw *ResponseWriter) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func logWrap(bl logger.Logger, hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := bl.Fork()
		if logLevel := r.Header.Get("X-Log-Request"); logLevel != "" {
			lvl, err := logger.ParseLevel(logLevel)
			if err != nil {
				l.NoRepublish().Errorf("Invalid requested log level %s", logLevel)
			} else {
				l = l.Trace(lvl)
			}
		}
		if logToken := r.Header.Get("X-Log-Token"); logToken != "" {
			l.NoRepublish().Errorf("Log token: %s", logToken)
		}
		start := time.Now()
		path := r.URL.Path
		raw := r.URL.RawQuery
		rw := &ResponseWriter{w, l, 0}
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			res := &models.Error{
				Code: 500,
				Type: r.Method,
				Key:  r.URL.Path,
			}
			res.Errorf("Panic recovered: %v", err)
			res.Errorf("Stack trace:")
			stack := bufio.NewScanner(bytes.NewReader(rt.Stack()))
			for stack.Scan() {
				res.Errorf("%s", stack.Text())
			}
			rw.Errorf("%s", res)
			JsonResponse(rw, res.Code, res)
		}()
		hf(rw, r)
		latency := time.Now().Sub(start)
		method := r.Method
		statusCode := rw.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		l.NoRepublish().Debugf("API: st: %d lt: %13v m: %s %s",
			statusCode,
			latency,
			method,
			path,
		)
	}
}

// Mux provides a logger wrapper http.ServeMux
// to provide integration with the logger system.
type Mux struct {
	logger.Logger
	*http.ServeMux
}

func nf(w http.ResponseWriter, r *http.Request) {
	res := &models.Error{
		Code: http.StatusNotFound,
		Type: r.Method,
		Key:  r.URL.Path,
	}
	JsonResponse(w, res.Code, res)
}

// Handle Map by request type
func (m *Mux) HandleMap(path string, mh map[string]http.HandlerFunc) {
	vh := func(w http.ResponseWriter, r *http.Request) {
		h, ok := mh[r.Method]
		if !ok {
			nf(w, r)
		} else {
			h(w, r)
		}
	}
	m.ServeMux.Handle(path, logWrap(m, vh))
}

// Handle registers a handler at the path on the provided Mux.
func (m *Mux) Handle(path string, h http.HandlerFunc) {
	vh := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			nf(w, r)
		} else {
			h(w, r)
		}
	}
	m.ServeMux.Handle(path, logWrap(m, vh))
}

// New creates a new Mux structure with the provided logger.
func New(bl logger.Logger) *Mux {
	res := &Mux{
		Logger:   bl,
		ServeMux: http.NewServeMux(),
	}
	res.Handle("/", nf)
	return res
}

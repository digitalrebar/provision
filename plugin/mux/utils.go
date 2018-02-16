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

func JsonResponse(w http.ResponseWriter, code int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	return enc.Encode(obj)
}

func TestContentType(r *http.Request, ct string) (string, bool) {
	rct := r.Header.Get("Content-Type")
	return rct, strings.Contains(
		strings.ToUpper(rct),
		strings.ToUpper(ct))
}

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

type ResponseWriter struct {
	http.ResponseWriter
	logger.Logger
	status int
}

func (rw *ResponseWriter) Status() int {
	return rw.status
}

func (rw *ResponseWriter) Write(buf []byte) (int, error) {
	if rw.status == 0 {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(buf)
}

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
				l.Errorf("Invalid requested log level %s", logLevel)
			} else {
				l = l.Trace(lvl)
			}
		}
		if logToken := r.Header.Get("X-Log-Token"); logToken != "" {
			l.Errorf("Log token: %s", logToken)
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
		l.Debugf("API: st: %d lt: %13v m: %s %s",
			statusCode,
			latency,
			method,
			path,
		)
	}
}

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

func New(bl logger.Logger) *Mux {
	res := &Mux{
		Logger:   bl,
		ServeMux: http.NewServeMux(),
	}
	res.Handle("/", nf)
	return res
}

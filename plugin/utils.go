package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/digitalrebar/logger"
	"github.com/digitalrebar/provision/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func testContentType(c *gin.Context, ct string) bool {
	ct = strings.ToUpper(ct)
	test := strings.ToUpper(c.ContentType())

	return strings.Contains(test, ct)
}

func assureContentType(c *gin.Context, ct string) bool {
	if testContentType(c, ct) {
		return true
	}
	err := &models.Error{Type: c.Request.Method, Code: http.StatusBadRequest}
	err.Errorf("Invalid content type: %s", c.ContentType())
	c.JSON(err.Code, err)
	return false
}

func assureDecode(c *gin.Context, val interface{}) bool {
	if !assureContentType(c, "application/json") {
		return false
	}
	if c.Request.ContentLength == 0 {
		val = nil
		return true
	}
	marshalErr := binding.JSON.Bind(c.Request, &val)
	if marshalErr == nil {
		return true
	}
	err := &models.Error{Type: c.Request.Method, Code: http.StatusBadRequest}
	err.AddError(marshalErr)
	c.JSON(err.Code, err)
	return false
}

func post(path string, indata interface{}) ([]byte, error) {
	if data, err := json.Marshal(indata); err != nil {
		return nil, err
	} else {
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
			berr := models.Error{}
			err := json.Unmarshal(b, &berr)
			if err != nil {
				return nil, e
			}
			return nil, &berr
		}

		return b, nil
	}
}

func newGinServer(bl logger.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	mgmtApi := gin.New()
	mgmtApi.Use(func(c *gin.Context) {
		l := bl.Fork()
		if logLevel := c.GetHeader("X-Log-Request"); logLevel != "" {
			lvl, err := logger.ParseLevel(logLevel)
			if err != nil {
				l.Errorf("Invalid requested log level %s", logLevel)
			} else {
				l = l.Trace(lvl)
			}
		}
		if logToken := c.GetHeader("X-Log-Token"); logToken != "" {
			l.Errorf("Log token: %s", logToken)
		}
		c.Set("logger", l)
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		latency := time.Now().Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		l.Debugf("API: st: %d lt: %13v ip: %15s m: %s %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	})
	mgmtApi.Use(gin.Recovery())

	return mgmtApi
}

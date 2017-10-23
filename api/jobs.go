package api

// Come back to processJobs later

import (
	"bytes"
	"fmt"
	"io"

	"github.com/digitalrebar/provision/models"
)

var exitOnFailure = false

func (c *Client) AppendJobLog(j *models.Job, buf io.Reader) error {
	uri := c.UrlFor("jobs", j.Key(), "log")
	req, err := c.Request("PUT", uri, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return Unmarshal(resp, nil)
}

func (c *Client) JobLogString(j *models.Job, s string, items ...interface{}) error {
	uri := c.UrlFor("jobs", j.Key(), "log")
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, s, items...)
	req, err := c.Request("PUT", uri, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return Unmarshal(resp, nil)
}

func (c *Client) JobLog(j *models.Job, dst io.Writer) error {
	uri := c.UrlFor("jobs", j.Key(), "log")
	req, err := c.Request("GET", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octest-stream")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return Unmarshal(resp, dst)
}

func (c *Client) JobActions(j *models.Job) ([]*models.JobAction, error) {
	uri := c.UrlFor("jobs", j.Key(), "actions")
	res := []*models.JobAction{}
	return res, c.DoJSON("GET", uri, nil, &res)
}

/*

func (c *Client) RunJobCommand(j *models.Job, action *models.JobAction) (failed, incomplete, reboot bool) {
	logBody, logWriter := io.Pipe()
	logger := bufio.NewWriter(logWriter)

	c.JobLogString(j, "")
}
*/

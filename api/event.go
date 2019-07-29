package api

import "github.com/digitalrebar/provision/v4/models"

func (c *Client) PostEvent(evt *models.Event) error {
	return c.Req().Post(evt).UrlFor("events").Do(nil)
}

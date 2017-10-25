package api

import "github.com/digitalrebar/provision/models"

func (c *Client) PostEvent(evt *models.Event) error {
	return c.Req().Post(evt).UrlFor("events").Do(nil)
}

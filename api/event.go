package api

import "github.com/digitalrebar/provision/models"

func (c *Client) PostEvent(evt *models.Event) error {
	return c.DoJSON("POST", c.UrlFor("events"), evt, nil)
}

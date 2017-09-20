package api

import (
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/digitalrebar/provision/models"
	"github.com/gorilla/websocket"
)

func (c *Client) ws() (*websocket.Conn, error) {
	ep := c.UrlFor("ws")
	ep.Scheme = "wss"
	dialer := &websocket.Dialer{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	header := http.Header{}
	header.Set("Authorization", "Bearer "+c.Token())
	res, _, err := dialer.Dial(ep.String(), header)
	return res, err
}

type RecievedEvent struct {
	E   models.Event
	Err error
}
type EventStream struct {
	conn   *websocket.Conn
	Events <-chan RecievedEvent
}

func (c *Client) Events() (*EventStream, error) {
	conn, err := c.ws()
	if err != nil {
		return nil, err
	}
	events := make(chan RecievedEvent)
	res := &EventStream{
		conn:   conn,
		Events: events,
	}
	go func(conn *websocket.Conn, events chan RecievedEvent) {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
				events <- RecievedEvent{Err: err}
				close(events)
				return
			}
			evt := RecievedEvent{}
			evt.Err = json.Unmarshal(msg, &evt.E)
			if err != nil {
				continue
			}
			events <- evt
		}
	}(conn, events)
	return res, nil
}

func (es *EventStream) Close() error {
	return es.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (es *EventStream) Register(events ...string) error {
	for _, evt := range events {
		if err := es.conn.WriteMessage(websocket.TextMessage, []byte("register "+evt)); err != nil {
			return err
		}
	}
	return nil
}

func (es *EventStream) Deregister(events ...string) error {
	for _, evt := range events {
		if err := es.conn.WriteMessage(websocket.TextMessage, []byte("deregister "+evt)); err != nil {
			return err
		}
	}
	return nil
}

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

// RecievedEvent contains an event recieved from the digitalrebar
// provision server along with any errors that occurred while
// recieving the event.
type RecievedEvent struct {
	E   models.Event
	Err error
}

// EventStream recieves events from the digitalrebar provider.  You can read recieved events by reading from its Events channel.
type EventStream struct {
	conn   *websocket.Conn
	Events <-chan RecievedEvent
}

// Events creates a new EventStream from the client.
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

// Close closes down the EventStream.  You should drain the Events
// until you read a RecievedEvent that has an empty E and a non-nil
// Err
func (es *EventStream) Close() error {
	return es.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

// Register directs the EventStream to subscribe to Events from the digital rebar provisioner.
//
// Event subscriptions consist of a string with the following format:
//
//    type.action.key
//
// type is the object type that you want to listen for events about.
// * means to listen for events about all object types.
//
// action is the action that caused the event to be created.  * means
// to listen for all actions.
//
// key is the unique identifier of the object to listen for.  * means
// to listen for events from all objects
func (es *EventStream) Register(events ...string) error {
	for _, evt := range events {
		if err := es.conn.WriteMessage(websocket.TextMessage, []byte("register "+evt)); err != nil {
			return err
		}
	}
	return nil
}

// Deregister directs the EventStream to unsubscribe from Events from
// the digitalrebar provisioner.  It takes the same parameters as
// Register.
func (es *EventStream) Deregister(events ...string) error {
	for _, evt := range events {
		if err := es.conn.WriteMessage(websocket.TextMessage, []byte("deregister "+evt)); err != nil {
			return err
		}
	}
	return nil
}

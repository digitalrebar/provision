package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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

func (r *RecievedEvent) matches(registration string) bool {
	tak := strings.SplitN(registration, ".", 3)
	if len(tak) != 3 {
		return false
	}
	return (tak[0] == r.E.Type || tak[0] == "*") &&
		(tak[1] == r.E.Action || tak[1] == "*") &&
		(tak[2] == r.E.Key || tak[2] == "*")
}

// EventStream recieves events from the digitalrebar provider.  You can read recieved events by reading from its Events channel.
type EventStream struct {
	handleId      int64
	conn          *websocket.Conn
	subscriptions map[string][]int64
	recievers     map[int64]chan RecievedEvent
	mux           *sync.Mutex
}

func (es *EventStream) processEvents(running chan struct{}) {
	close(running)
	for {
		log.Printf("waiting on eventstream msg")
		_, msg, err := es.conn.ReadMessage()
		log.Printf("Recieved msg: %s", string(msg))
		if err != nil {
			log.Printf("Recieved msg has error: %v", err)
			log.Printf("Closing down channels")
			es.conn.Close()
			es.mux.Lock()
			for _, reciever := range es.recievers {
				reciever <- RecievedEvent{Err: err}
				close(reciever)
			}
			es.mux.Unlock()
			return
		}
		evt := RecievedEvent{}
		evt.Err = json.Unmarshal(msg, &evt.E)
		toSend := map[int64]chan RecievedEvent{}
		log.Printf("Locking es mux in process loop")
		es.mux.Lock()
		for reg, handles := range es.subscriptions {
			if !evt.matches(reg) {
				log.Printf("Evt %v does not match %s, skipping", evt.E, reg)
				continue
			}
			log.Printf("Evt %v matches %s, queing it up to send", evt.E, reg)
			for _, i := range handles {
				if toSend[i] == nil {
					toSend[i] = es.recievers[i]
				}
			}
		}
		es.mux.Unlock()
		log.Printf("es mux in process loop unlocked")
		go func(ts map[int64]chan RecievedEvent, evt RecievedEvent) {
			for i := range ts {
				log.Printf("Sending evt %v to handle %d", evt.E, i)
				ts[i] <- evt
			}
		}(toSend, evt)
	}
}

// Events creates a new EventStream from the client.
func (c *Client) Events() (*EventStream, error) {
	conn, err := c.ws()
	if err != nil {
		return nil, err
	}
	res := &EventStream{
		conn:          conn,
		subscriptions: map[string][]int64{},
		recievers:     map[int64]chan RecievedEvent{},
		mux:           &sync.Mutex{},
	}
	running := make(chan struct{})
	go res.processEvents(running)
	<-running
	return res, nil
}

// Close closes down the EventStream.  You should drain the Events
// until you read a RecievedEvent that has an empty E and a non-nil
// Err
func (es *EventStream) Close() error {
	return es.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (es *EventStream) subscribe(handle int64, events ...string) error {
	if es.recievers[handle] == nil {
		return fmt.Errorf("No such handle %d", handle)
	}
	for _, evt := range events {
		handles := es.subscriptions[evt]
		if handles == nil {
			handles = []int64{}
		}
		idx := sort.Search(len(handles), func(i int) bool { return handles[i] >= handle })
		if idx == len(handles) {
			handles = append(handles, handle)
		} else if handles[idx] == handle {
			continue
		} else {
			handles = append(handles, 0)
			copy(handles[idx+1:], handles[idx:])
			handles[idx] = handle
		}
		if es.subscriptions[evt] == nil {
			if err := es.conn.WriteMessage(websocket.TextMessage, []byte("register "+evt)); err != nil {
				return err
			}
		}
		es.subscriptions[evt] = handles
	}
	return nil
}

func (es *EventStream) Subscribe(handle int64, events ...string) error {
	es.mux.Lock()
	defer es.mux.Unlock()
	return es.subscribe(handle, events...)
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
func (es *EventStream) Register(events ...string) (int64, <-chan RecievedEvent, error) {
	newID := atomic.AddInt64(&es.handleId, 1)
	es.mux.Lock()
	defer es.mux.Unlock()
	ch := make(chan RecievedEvent)
	es.recievers[newID] = ch
	return newID, ch, es.subscribe(newID, events...)
}

// Deregister directs the EventStream to unsubscribe from Events from
// the digitalrebar provisioner.  It takes the same parameters as
// Register.
func (es *EventStream) Deregister(handle int64) error {
	es.mux.Lock()
	defer es.mux.Unlock()
	ch := es.recievers[handle]
	if ch == nil {
		return fmt.Errorf("No such handle %d", handle)
	}
	for evt, handles := range es.subscriptions {
		idx := sort.Search(len(handles), func(i int) bool { return handles[i] >= handle })
		if idx == len(handles) || handles[idx] != handle {
			continue
		} else if idx != len(handles)-1 {
			copy(handles[idx:], handles[idx+1:])
		}
		handles = handles[:len(handles)-1]
		es.subscriptions[evt] = handles
		if len(handles) == 0 {
			es.conn.WriteMessage(websocket.TextMessage, []byte("deregister "+evt))
		}
	}
	delete(es.recievers, handle)
	close(ch)
	return nil
}

// WaitFor waits for an item to match test.  It subscribes to an
// EventStream that watches all update and save envents for the object
// in question, and returns a string indicating whether the match
// succeeded, failed, or timed out.
//
// The API for this function is subject to refactoring and change, and
// should not be considered to be stable yet.
func (es *EventStream) WaitFor(item models.Model, test func(interface{}) (bool, error), timeout int64) (string, error) {
	prefix := item.Prefix()
	id := item.Key()
	interrupt := make(chan os.Signal, 1)
	evts := []string{prefix + ".update." + id, prefix + ".save." + id}
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	handle, ch, err := es.Register(evts...)
	defer es.Deregister(handle)
	if err != nil {
		return "", err
	}
	timer := time.NewTimer(time.Second * time.Duration(timeout))
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()
	for {
		found, err := test(item)
		if found && err == nil {
			return "complete", nil
		}
		if err != nil {
			return fmt.Sprintf("test: %v", err), err
		}
		select {
		case evt := <-ch:
			if evt.Err != nil {
				return fmt.Sprintf("read: %v", err), err
			}
			item, err = evt.E.Model()
			if err != nil {
				return fmt.Sprintf("read: %v", err), err
			}
		case <-interrupt:
			return "interrupt", nil
		case <-timer.C:
			return "timeout", nil
		}
	}
}

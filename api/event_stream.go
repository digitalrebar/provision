package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision/models"
	"github.com/gorilla/websocket"
)

// TestItem creates a test function to see if a value in the
// passed interface is true.
func TestItem(field, value string) func(interface{}) (bool, error) {
	return func(ref interface{}) (bool, error) {
		var err error
		fields := map[string]interface{}{}
		if err := utils.Remarshal(ref, &fields); err != nil {
			return false, err
		}
		matched := false
		if d, ok := fields[field]; ok {
			switch v := d.(type) {
			case bool:
				var bval bool
				bval, err = strconv.ParseBool(value)
				if err == nil {
					if v == bval {
						matched = true
					}
				}
			case string:
				if v == value {
					matched = true
				}
			case int:
				var ival int64
				ival, err = strconv.ParseInt(value, 10, 64)
				if err == nil {
					if int(ival) == v {
						matched = true
					}
				}
			default:
				err = fmt.Errorf("Unsupported field type: %T\n", d)
			}
		}
		return matched, err
	}
}

func (c *Client) ws() (*websocket.Conn, error) {
	ep, err := c.UrlFor("ws")
	if err != nil {
		return nil, err
	}
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
	client        *Client
	handleId      int64
	conn          *websocket.Conn
	subscriptions map[string][]int64
	recievers     map[int64]chan RecievedEvent
	mux           *sync.Mutex
}

func (es *EventStream) processEvents(running chan struct{}) {
	close(running)
	for {
		_, msg, err := es.conn.ReadMessage()
		if err != nil {
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
		es.mux.Lock()
		for reg, handles := range es.subscriptions {
			if !evt.matches(reg) {
				continue
			}
			for _, i := range handles {
				if toSend[i] == nil {
					toSend[i] = es.recievers[i]
				}
			}
		}
		es.mux.Unlock()
		go func(ts map[int64]chan RecievedEvent, evt RecievedEvent) {
			for i := range ts {
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
		client:        c,
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
func (es *EventStream) WaitFor(
	item models.Model,
	test func(interface{}) (bool, error),
	timeout time.Duration) (string, error) {
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
	timer := time.NewTimer(timeout)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()
	for {
		if err := es.client.FillModel(item, id); err != nil {
			return fmt.Sprintf("fill: %v", err), err
		}
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

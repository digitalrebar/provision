package api

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision/models"
	"github.com/gorilla/websocket"
)

// TestFunc is a function type used to test if an item matches some
// sort of condition.
type TestFunc func(interface{}) (bool, error)

// AndItems retuens a TestFunc that returns true if all the passed-in
// TestFuncs also return true.
func AndItems(fs ...TestFunc) TestFunc {
	return func(ref interface{}) (bool, error) {
		for _, f := range fs {
			if b, e := f(ref); e != nil {
				return false, e
			} else if !b {
				return false, nil
			}
		}
		return true, nil
	}
}

// OrItems returns a TestFunc that returns true if any of the
// passed-in TestFuncs return true.
func OrItems(fs ...TestFunc) TestFunc {
	return func(ref interface{}) (bool, error) {
		for _, f := range fs {
			if b, e := f(ref); e != nil {
				return false, e
			} else if b {
				return true, e
			}
		}
		return false, nil
	}
}

// NotItems returns a TestFunc that returns the opposite of the
// passed-in TestFunc.
func NotItem(f TestFunc) TestFunc {
	return func(ref interface{}) (bool, error) {
		b, e := f(ref)
		return !b, e
	}
}

// EqualItem creates a test function to see if a value in
// the passed interface is equal
func EqualItem(field string, value interface{}) TestFunc {
	// Convert the value to a json remarshaled value
	var nv interface{}
	if err := utils.Remarshal(value, &nv); err != nil {
		return func(ref interface{}) (bool, error) {
			return false, err
		}
	}

	return func(ref interface{}) (bool, error) {
		fields := map[string]interface{}{}
		if err := utils.Remarshal(ref, &fields); err != nil {
			return false, err
		}
		if d, ok := fields[field]; ok {
			return reflect.DeepEqual(d, nv), nil
		}
		return false, nil
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
	// If we have a token use it, otherwise basic auth
	if c.Token() != "" {
		header.Set("Authorization", "Bearer "+c.Token())
	} else {
		basicAuth := base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.password))
		header.Set("Authorization", "Basic "+basicAuth)
	}
	res, _, err := dialer.Dial(ep.String(), header)
	return res, err
}

// RecievedEvent contains an event received from the digitalrebar
// provision server along with any errors that occurred while
// receiving the event.
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

// EventStream receives events from the digitalrebar provider.  You
// can read received events by reading from its Events channel.
type EventStream struct {
	client        *Client
	handleId      int64
	conn          *websocket.Conn
	subscriptions map[string][]int64
	receivers     map[int64]chan RecievedEvent
	mux           *sync.Mutex
	rchan         chan RecievedEvent
}

func (es *EventStream) processEvents(running chan struct{}) {
	close(running)
	for {
		_, msg, err := es.conn.NextReader()
		if err != nil {
			es.conn.Close()
			es.mux.Lock()
			for h, receiver := range es.receivers {
				receiver <- RecievedEvent{Err: err}
				close(receiver)
				es.receivers[h] = nil
			}
			es.mux.Unlock()
			return
		}
		evt := RecievedEvent{}
		evt.Err = json.NewDecoder(msg).Decode(&evt.E)
		toSend := map[int64]chan RecievedEvent{}
		es.mux.Lock()
		for reg, handles := range es.subscriptions {
			if !evt.matches(reg) {
				continue
			}
			for _, i := range handles {
				if toSend[i] == nil {
					toSend[i] = es.receivers[i]
				}
			}
		}
		for i := range toSend {
			select {
			case toSend[i] <- evt:
			default:
				fmt.Printf("Failed to send an event\n")
			}
		}
		es.mux.Unlock()
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
		receivers:     map[int64]chan RecievedEvent{},
		mux:           &sync.Mutex{},
	}
	newID := atomic.AddInt64(&res.handleId, 1)
	res.rchan = make(chan RecievedEvent, 100)
	res.subscriptions["websocket.*.*"] = []int64{newID}
	res.receivers[newID] = res.rchan
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

func (es *EventStream) subscribe(handle int64, events ...string) (int, error) {
	if es.receivers[handle] == nil {
		return 0, fmt.Errorf("No such handle %d", handle)
	}
	count := 0
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
				return count, err
			}
			count += 1
		}
		es.subscriptions[evt] = handles
	}
	return count, nil
}

// Subscribe directs an EventStream to listen for matching events.
func (es *EventStream) Subscribe(handle int64, events ...string) error {
	es.mux.Lock()
	count, err := es.subscribe(handle, events...)
	es.mux.Unlock()
	if err == nil {
		// Really wait should be for my specific one, but ...
		// Multi-threaded apps will have issues.
		for i := 0; i < count; i++ {
			<-es.rchan
		}
	}
	return err
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
	ch := make(chan RecievedEvent, 100)
	es.receivers[newID] = ch
	count, err := es.subscribe(newID, events...)
	es.mux.Unlock()
	if err == nil {
		// Really wait should be for my specific one, but ...
		// Multi-threaded apps will have issues.
		for i := 0; i < count; i++ {
			<-es.rchan
		}
	}
	return newID, ch, err
}

func (es *EventStream) deregister(handle int64) (int, error) {
	ch, ok := es.receivers[handle]
	if !ok {
		return 0, fmt.Errorf("No such handle %d", handle)
	}
	count := 0
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
			count += 1
			es.conn.WriteMessage(websocket.TextMessage, []byte("deregister "+evt))
			delete(es.subscriptions, evt)
		}
	}
	delete(es.receivers, handle)
	if ch != nil {
		close(ch)
	}
	return count, nil
}

// Deregister directs the EventStream to unsubscribe from Events from
// the digitalrebar provisioner.  It takes the same parameters as
// Register.
func (es *EventStream) Deregister(handle int64) error {
	es.mux.Lock()
	count, err := es.deregister(handle)
	es.mux.Unlock()
	// Really wait should be for my specific one, but ...
	// Multi-threaded apps will have issues.
	for i := 0; i < count; i++ {
		<-es.rchan
	}
	return err
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
	test TestFunc,
	timeout time.Duration) (string, error) {
	// Make some basic vars
	prefix := item.Prefix()
	id := item.Key()
	interrupt := make(chan os.Signal, 1)
	evts := []string{prefix + ".update." + id, prefix + ".save." + id}

	// Handle interrupt signal while selecting
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Reset(os.Interrupt)

	// Register for events
	handle, ch, err := es.Register(evts...)
	if err != nil {
		return "", err
	}
	defer es.Deregister(handle)

	// Setup the timer
	timer := time.NewTimer(timeout)
	defer func() {
		if timer != nil && !timer.Stop() {
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
			if err != nil {
				return fmt.Sprintf("read: %v", err), err
			}
		case <-interrupt:
			return "interrupt", nil
		case <-timer.C:
			timer.Stop()
			timer = nil
			return "timeout", nil
		}
	}
}

package api

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision/models"
)

func testItem(ref interface{}, field, value string) (bool, error) {
	var err error
	fields := map[string]interface{}{}
	if err := utils.Remarshal(ref, fields); err != nil {
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

func (c *Client) WaitFor(item models.Model, field, value string, timeout int64) (string, error) {
	found, err := testItem(item, field, value)
	if err != nil {
		return fmt.Sprintf("read: %v", err), err
	}
	if found {
		return "complete", nil
	}
	prefix := item.Prefix()
	id := item.Key()
	interrupt := make(chan os.Signal, 1)
	evts := []string{prefix + ".update." + id, prefix + ".save." + id}
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	es, err := c.Events()
	if err != nil {
		return "", err
	}

	if err := es.Register(evts...); err != nil {
		es.Close()
		return "", err
	}
	timer := time.NewTimer(time.Second * time.Duration(timeout))
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
		es.Deregister(evts...)
		es.Close()
	}()
	for {
		select {
		case evt := <-es.Events:
			if evt.Err != nil {
				return fmt.Sprintf("read: %v", err), err
			}
			obj, err := evt.E.Model()
			if err != nil {
				return fmt.Sprintf("read: %v", err), err
			}
			found, err := testItem(obj, field, value)
			if found && err == nil {
				return "complete", nil
			}
			if err != nil {
				return fmt.Sprintf("test: %v", err), err
			}
		case <-interrupt:
			return "interrupt", nil
		case <-timer.C:
			return "timeout", nil
		}
	}
}

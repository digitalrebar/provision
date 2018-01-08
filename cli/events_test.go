package cli

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/digitalrebar/provision/models"
)

func TestEventsCli(t *testing.T) {
	var eventsIntroString = "DigitalRebar Provision Event Commands\n"
	var eventsPostNoArgString = "Error: drpcli events post [- | JSON or YAML Event] [flags] requires 1 argument\n"
	var eventsPostTooManyArgsString = "Error: drpcli events post [- | JSON or YAML Event] [flags] requires 1 argument\n"
	var eventsPostBadJsonString = "Error: Invalid event: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n\n"
	var eventsPostBadJson1String = "Error: Invalid event: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Event\n\n\n"

	event := &models.Event{Time: time.Now(), Type: "events", Action: "post", Key: "test", Object: "String of Data"}
	jsonBytes, _ := json.Marshal(event)
	jsonString := string(jsonBytes)

	CliTest{true, false, []string{"events"}, noStdinString, eventsIntroString, noErrorString, ""}.run(t)
	CliTest{true, true, []string{"events", "post"}, noStdinString, noContentString, eventsPostNoArgString, ""}.run(t)
	CliTest{true, true, []string{"events", "post", "e1", "e2"}, noStdinString, noContentString, eventsPostTooManyArgsString, ""}.run(t)
	CliTest{false, true, []string{"events", "post", "{sasdg"}, noStdinString, noContentString, eventsPostBadJsonString, ""}.run(t)
	CliTest{false, true, []string{"events", "post", "\"e1\""}, noStdinString, noContentString, eventsPostBadJson1String, ""}.run(t)
	CliTest{false, false, []string{"events", "post", jsonString}, noStdinString, noContentString, noErrorString, ""}.run(t)
	CliTest{false, false, []string{"events", "post", "-"}, jsonString, noContentString, noErrorString, ""}.run(t)

}

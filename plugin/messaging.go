package midlayer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/digitalrebar/provision/backend"
)

type PluginClient struct {
	plugin   string
	cmd      *exec.Cmd
	stderr   io.ReadCloser
	stdout   io.ReadCloser
	stdin    io.WriteCloser
	finished chan bool
	logger   *log.Logger
	lock     sync.Mutex
	nextId   int
	pending  map[int]*PluginClientRequest
}

// Id of request, and JSON blob
type PluginClientRequest struct {
	Id     int
	Action string
	Data   interface{}

	caller chan string
}

type PluginClientReply struct {
	Id   int
	Code int
	Data interface{}
}

func (pc *PluginClient) ReadLog() {
	// read command's stderr line by line - for logging
	in := bufio.NewScanner(pc.stderr)
	for in.Scan() {
		pc.logger.Printf("Plugin " + pc.plugin + ": " + in.Text()) // write each line to your log, or anything you need
	}
	if err := in.Err(); err != nil {
		pc.logger.Printf("Plugin %s: error: %s", pc.plugin, err)
	}
	pc.finished <- true
}

func (pc *PluginClient) ReadReply() {
	// read command's stdout line by line - for replies
	in := bufio.NewScanner(pc.stdout)
	for in.Scan() {
		jsonString := in.Text()

		var resp PluginClientReply
		err := json.Unmarshal([]byte(jsonString), &resp)
		if err != nil {
			pc.logger.Printf("Failed to process: %v\n", err)
			continue
		}

		req, ok := pc.pending[resp.Id]
		if !ok {
			pc.logger.Printf("Failed to find request for: %v\n", resp.Id)
			continue
		}

		req.caller <- jsonString

		pc.lock.Lock()
		delete(pc.pending, resp.Id)
		pc.lock.Unlock()
	}
	if err := in.Err(); err != nil {
		pc.logger.Printf("Reply %s: error: %s", pc.plugin, err)
	}
	pc.finished <- true
}

func (pc *PluginClient) writeRequest(action string, data interface{}) (chan string, error) {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	mychan := make(chan string)
	id := pc.nextId
	pc.pending[id] = &PluginClientRequest{Id: id, Action: action, Data: data, caller: mychan}
	pc.nextId += 1

	if bytes, err := json.Marshal(pc.pending[id]); err != nil {
		delete(pc.pending, id)
		return mychan, nil
	} else {
		n, err := pc.stdin.Write(bytes)
		if err != nil {
			return mychan, err
		}
		if n != len(bytes) {
			return mychan, fmt.Errorf("Failed to write all bytes: %d (%d)\n", len(bytes), n)
		}
		n, err = pc.stdin.Write([]byte("\n"))
		if err != nil {
			return mychan, err
		}
	}

	return mychan, nil
}

func (pc *PluginClient) Config(params map[string]interface{}) error {
	if mychan, err := pc.writeRequest("Config", params); err != nil {
		return err
	} else {
		s := <-mychan
		pc.logger.Printf("GREG: Config Reply: %s\n", s)
	}
	return nil
}

func (pc *PluginClient) Publish(e *backend.Event) error {
	if mychan, err := pc.writeRequest("Publish", e); err != nil {
		return err
	} else {
		answer := <-mychan
		pc.logger.Printf("GREG: Publish Reply: %s\n", answer)
	}
	return nil
}

func (pc *PluginClient) Action(a *MachineAction) error {
	if mychan, err := pc.writeRequest("Action", a); err != nil {
		return err
	} else {
		answer := <-mychan
		pc.logger.Printf("GREG: Action Reply: %s\n", answer)
	}
	return nil
}

func (pc *PluginClient) Stop() error {
	// Close stdin / writer.  To close, the program.
	pc.logger.Printf("GREG: Stopping program by closing STDIN\n")
	pc.stdin.Close()

	// Wait for reader to exit
	pc.logger.Printf("GREG: Waiting for log reader to finish\n")
	<-pc.finished
	pc.logger.Printf("GREG: Waiting for reply reader to finish\n")
	<-pc.finished

	// Wait for exit
	pc.logger.Printf("GREG: Wait for command to exit\n")
	pc.cmd.Wait()
	return nil
}

func NewPluginClient(plugin string, logger *log.Logger, apiPort int, path string, params map[string]interface{}) (answer *PluginClient, theErr error) {
	answer = &PluginClient{plugin: plugin, logger: logger, pending: make(map[int]*PluginClientRequest, 0)}

	answer.cmd = exec.Command(path, "listen")
	// Setup env vars to run drpcli - auth should be parameters.
	env := os.Environ()
	env = append(env, fmt.Sprintf("RS_ENDPOINT=https://127.0.0.1:%d", apiPort))
	answer.cmd.Env = env

	var err2 error
	answer.stderr, err2 = answer.cmd.StderrPipe()
	if err2 != nil {
		return nil, err2
	}
	answer.stdout, err2 = answer.cmd.StdoutPipe()
	if err2 != nil {
		return nil, err2
	}
	answer.stdin, err2 = answer.cmd.StdinPipe()
	if err2 != nil {
		return nil, err2
	}

	answer.finished = make(chan bool, 2)
	go answer.ReadLog()
	go answer.ReadReply()

	answer.cmd.Start()

	answer.logger.Printf("GREG: Calling plugin.config %v\n", params)
	terr := answer.Config(params)
	answer.logger.Printf("GREG: returned plugin.config %v\n", terr)
	if terr != nil {
		answer.Stop()
		theErr = terr
		return
	}
	return
}

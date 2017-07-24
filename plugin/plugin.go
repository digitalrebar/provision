package plugin

/*
 * This is used by plugins to define their base App.
 */

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/digitalrebar/provision/backend"
	"github.com/spf13/cobra"
)

type PluginConfig interface {
	Config(map[string]interface{}) *backend.Error
}

type PluginPublisher interface {
	Publish(*backend.Event) *backend.Error
}

type PluginActor interface {
	Action(*MachineAction) *backend.Error
}

var (
	App = &cobra.Command{
		Use:   "replaceme",
		Short: "Replace ME!",
	}
	debug = false
)

func Log(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func InitApp(use, short, version string, def *PluginProvider, pc PluginConfig) {
	App.Use = use
	App.Short = short

	App.PersistentFlags().BoolVarP(&debug,
		"debug", "d", false,
		"Whether the CLI should run in debug mode")

	App.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Digital Rebar Provision CLI Command Version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version: %v\n", version)
			return nil
		},
	})
	App.AddCommand(&cobra.Command{
		Use:   "autocomplete <filename>",
		Short: "Digital Rebar Provision CLI Command Bash AutoCompletion File",
		Long:  "Generate a bash autocomplete file as <filename>.\nPlace the generated file in /etc/bash_completion.d or /usr/local/etc/bash_completion.d.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1  argument", cmd.UseLine())
			}
			App.GenBashCompletionFile(args[0])
			return nil
		},
	})
	App.AddCommand(&cobra.Command{
		Use:   "define",
		Short: "Digital Rebar Provision CLI Command Define",
		RunE: func(cmd *cobra.Command, args []string) error {
			prettyPrint(def)
			return nil
		},
	})
	App.AddCommand(&cobra.Command{
		Use:   "listen",
		Short: "Digital Rebar Provision CLI Command Listen",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(pc)
		},
	})
}

func prettyPrint(o interface{}) (err error) {
	var buf []byte
	buf, err = json.MarshalIndent(o, "", "  ")
	fmt.Println(string(buf))
	return nil
}

func Run(pc PluginConfig) error {
	var ops int64 = 0
	origStdOut := os.Stdout

	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		jsonString := in.Text()

		var req PluginClientRequest
		err := json.Unmarshal([]byte(jsonString), &req)
		if err != nil {
			Log("Failed to process: %v\n", err)
			continue
		}

		atomic.AddInt64(&ops, 1)
		go handleRequest(pc, &req, &ops, origStdOut)
	}
	if err := in.Err(); err != nil {
		Log("Plugin error: %s", err)
	}

	for a := atomic.LoadInt64(&ops); a > 0; {
		Log("Exiting ... waiting on %d go routines to leave\n", a)
		time.Sleep(time.Second * 1)
	}

	return nil
}

func handleRequest(pc PluginConfig, req *PluginClientRequest, ops *int64, origStdOut *os.File) {
	defer atomic.AddInt64(ops, -1)

	if req.Action == "Config" {
		resp := &PluginClientReply{Id: req.Id}

		var params map[string]interface{}
		if jerr := json.Unmarshal(req.Data, &params); jerr != nil {
			resp.Code = 400
			resp.Data, _ = json.Marshal(&backend.Error{Code: 400, Model: "plugin", Type: "Config", Messages: []string{fmt.Sprintf("Failed to unmarshal data: %v", jerr)}})
		} else {
			if err := pc.Config(params); err != nil {
				resp.Code = err.Code
				resp.Data, _ = json.Marshal(err)
			}
		}
		bytes, _ := json.Marshal(resp)
		fmt.Fprintln(origStdOut, string(bytes))
	} else if req.Action == "Action" {
		resp := &PluginClientReply{Id: req.Id}

		var actionInfo MachineAction
		if jerr := json.Unmarshal(req.Data, &actionInfo); jerr != nil {
			resp.Code = 400
			resp.Data, _ = json.Marshal(&backend.Error{Code: 400, Model: "plugin", Type: "Action", Messages: []string{fmt.Sprintf("Failed to unmarshal data: %v", jerr)}})
		} else {
			s, ok := pc.(PluginActor)
			if !ok {
				resp.Code = 400
				resp.Data, _ = json.Marshal(&backend.Error{Code: 400, Model: "plugin", Type: "Action", Messages: []string{"Plugin doesn't support Action"}})
			} else {
				if err := s.Action(&actionInfo); err != nil {
					resp.Code = err.Code
					resp.Data, _ = json.Marshal(err)
				}
			}
		}
		bytes, _ := json.Marshal(resp)
		fmt.Fprintln(origStdOut, string(bytes))
	} else if req.Action == "Publish" {
		resp := &PluginClientReply{Id: req.Id}

		var event backend.Event
		if jerr := json.Unmarshal(req.Data, &event); jerr != nil {
			resp.Code = 400
			resp.Data, _ = json.Marshal(&backend.Error{Code: 400, Model: "plugin", Type: "Publish", Messages: []string{fmt.Sprintf("Failed to unmarshal data: %v", jerr)}})
		} else {
			s, ok := pc.(PluginPublisher)
			if !ok {
				resp.Code = 400
				resp.Data, _ = json.Marshal(&backend.Error{Code: 400, Model: "plugin", Type: "Publish", Messages: []string{"Plugin doesn't support Publish"}})
			} else {
				if err := s.Publish(&event); err != nil {
					resp.Code = err.Code
					resp.Data, _ = json.Marshal(err)
				}
			}
		}
		bytes, _ := json.Marshal(resp)
		fmt.Fprintln(origStdOut, string(bytes))
	} else {
		resp := &PluginClientReply{Code: 400, Id: req.Id}
		resp.Data, _ = json.Marshal(&backend.Error{Code: 400, Model: "plugin", Type: "Publish", Messages: []string{"Plugin unknown command type"}})
		bytes, _ := json.Marshal(resp)
		fmt.Fprintln(origStdOut, string(bytes))
	}
}

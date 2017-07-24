package plugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/midlayer"
	"github.com/spf13/cobra"
)

type PluginConfig interface {
	Config(map[string]interface{}) *backend.Error
}

type PluginPublisher interface {
	Publish(*backend.Event) *backend.Error
}

type PluginActor interface {
	Action(*midlayer.MachineAction) *backend.Error
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

func InitApp(use, short, version string, def *midlayer.PluginProvider, pc PluginConfig) {
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
	// read command's stdin line by line - for logging
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		jsonString := in.Text()

		var req midlayer.PluginClientRequest
		err := json.Unmarshal([]byte(jsonString), &req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to process: %v\n", err)
			continue
		}

		if req.Action == "Config" {
			params := make(map[string]interface{}, 0)
			if req.Data != nil {
				params = req.Data.(map[string]interface{})
			}
			err := pc.Config(params)
			code := 0
			if err != nil {
				code = err.Code
			}

			resp := &midlayer.PluginClientReply{Code: code, Id: req.Id, Data: err}
			bytes, _ := json.Marshal(resp)
			fmt.Println(string(bytes))
		} else if req.Action == "Action" {
			fmt.Fprintf(os.Stderr, "GREG: Data type = %V\n", req.Data)
			fmt.Fprintf(os.Stderr, "GREG: Data type = %v\n", req.Data)
			fmt.Fprintf(os.Stderr, "GREG: Data type = %T\n", req.Data)
			fmt.Fprintf(os.Stderr, "GREG: Data type = %t\n", req.Data)
			actionInfo, ok := req.Data.(midlayer.MachineAction)
			if !ok {
				resp := &midlayer.PluginClientReply{Code: 400, Id: req.Id, Data: "Unknown data type"}
				bytes, _ := json.Marshal(resp)
				fmt.Println(string(bytes))
				continue
			}

			s, ok := pc.(PluginActor)
			if !ok {
				resp := &midlayer.PluginClientReply{Code: 400, Id: req.Id, Data: "Plugin doesn't support Action"}
				bytes, _ := json.Marshal(resp)
				fmt.Println(string(bytes))
				continue
			}

			err := s.Action(&actionInfo)
			code := 0
			if err != nil {
				code = err.Code
			}

			resp := &midlayer.PluginClientReply{Code: code, Id: req.Id, Data: err}
			bytes, _ := json.Marshal(resp)
			fmt.Println(string(bytes))
		} else if req.Action == "Publish" {
			event, ok := req.Data.(backend.Event)
			if !ok {
				resp := &midlayer.PluginClientReply{Code: 400, Id: req.Id, Data: "Unknown data type"}
				bytes, _ := json.Marshal(resp)
				fmt.Println(string(bytes))
				continue
			}

			s, ok := pc.(PluginPublisher)
			if !ok {
				resp := &midlayer.PluginClientReply{Code: 400, Id: req.Id, Data: "Plugin doesn't support Publish"}
				bytes, _ := json.Marshal(resp)
				fmt.Println(string(bytes))
				continue
			}

			err := s.Publish(&event)
			code := 0
			if err != nil {
				code = err.Code
			}

			resp := &midlayer.PluginClientReply{Code: code, Id: req.Id, Data: err}
			bytes, _ := json.Marshal(resp)
			fmt.Println(string(bytes))
		} else {
			resp := &midlayer.PluginClientReply{Code: 400, Id: req.Id, Data: "Unknown op"}
			bytes, _ := json.Marshal(resp)
			fmt.Println(string(bytes))
		}
	}
	if err := in.Err(); err != nil {
		fmt.Printf("Plugin error: %s", err)
	}

	return nil
}

package plugin

import (
	"encoding/json"
	"fmt"
	"net/rpc/jsonrpc"
	"os"

	"github.com/digitalrebar/provision/midlayer"
	"github.com/spf13/cobra"
)

type filePair struct {
	reader *os.File
	writer *os.File
}

func (fp *filePair) Read(p []byte) (int, error) {
	return fp.reader.Read(p)
}

func (fp *filePair) Write(p []byte) (int, error) {
	return fp.writer.Write(p)
}

func (fp *filePair) Close() error {
	fp.writer.Close()
	return fp.reader.Close()
}

var (
	App = &cobra.Command{
		Use:   "replaceme",
		Short: "Replace ME!",
	}
	debug = false
)

func InitApp(use, short, version string, def *midlayer.PluginProvider) {
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
			Run()
			// No return!
			return nil
		},
	})
}

func prettyPrint(o interface{}) (err error) {
	var buf []byte
	buf, err = json.MarshalIndent(o, "", "  ")
	fmt.Println(string(buf))
	return nil
}

func Run() {
	files := filePair{reader: os.Stdin, writer: os.Stdout}
	// GREG: errors and retry and exit?
	jsonrpc.ServeConn(&files)

}

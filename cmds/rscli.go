package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/rackn/rocket-skates/cli"
	apiclient "github.com/rackn/rocket-skates/client"
	"github.com/spf13/cobra"
)

func init() {
	if ep := os.Getenv("RS_ENDPOINT"); ep != "" {
		cli.Endpoint = ep
	}
	if kv := os.Getenv("RS_KEY"); kv != "" {
		key := strings.SplitN(kv, ":", 2)
		if len(key) < 2 {
			log.Fatal("RS_KEY does not contain a username:password pair!")
		}
		if key[0] == "" || key[1] == "" {
			log.Fatal("RS_KEY contains an invalid username:password pair!")
		}
		cli.Username = key[0]
		cli.Password = key[1]
	}
	cli.App.PersistentFlags().StringVarP(&cli.Endpoint,
		"endpoint", "E", cli.Endpoint,
		"The Rocket-Skates API endpoint to talk to")
	cli.App.PersistentFlags().StringVarP(&cli.Username,
		"username", "U", cli.Username,
		"Name of the Rocket-Skates user to talk to")
	cli.App.PersistentFlags().StringVarP(&cli.Password,
		"password", "P", cli.Password,
		"Password of the Rocket-Skates user")
	cli.App.PersistentFlags().BoolVarP(&cli.Debug,
		"debug", "d", false,
		"Whether the CLI should run in debug mode")
	cli.App.PersistentFlags().StringVarP(&cli.Format,
		"format", "F", "json",
		`The serialzation we expect for output.  Can be "json" or "yaml"`)
}

func main() {
	cli.App.PersistentPreRun = func(c *cobra.Command, a []string) {
		var err error
		cli.D("Talking to Rocket-Skates with %v (%v:%v)", cli.Endpoint, cli.Username, cli.Password)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		hc := &http.Client{Transport: tr}
		epURL, err := url.Parse(cli.Endpoint)
		if err != nil {
			log.Fatalf("Error handling endpoint %s: %v", cli.Endpoint, err)
		}
		transport := httptransport.NewWithClient(epURL.Host, "/api/v3", []string{epURL.Scheme}, hc)
		cli.Session = apiclient.New(transport, strfmt.Default)
		cli.BasicAuth = httptransport.BasicAuth(cli.Username, cli.Password)

		if err != nil {
			if c.Use != "version" {
				log.Fatalf("Could not connect to Rocket-Skates: %v\n", err.Error())
			}
		}
	}
	cli.App.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Rocket-Skates CLI Command Version",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Version: %v", cli.Version)
		},
	})
	cli.App.Execute()
}

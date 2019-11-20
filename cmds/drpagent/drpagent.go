package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/digitalrebar/provision/v4/agent"
	"github.com/digitalrebar/provision/v4/api"

	"github.com/digitalrebar/provision/v4/models"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	Endpoint   string `short:"E" long:"endpoint" env:"RS_ENDPOINT" description:"The URL of the dr-provision API endpoint to talk to"`
	Token      string `short:"T" long:"token" env:"RS_TOKEN" description:"The machine token to use for authentication."`
	MachineID  string `short:"m" long:"machine" env:"RS_UUID" description:"The UUID of the machine."`
	Context    string `short:"c" long:"context" env:"RS_CONTEXT" description:"The context that the agent should pay attention to"`
	StateDir   string `short:"s" long:"stateDir" env:"RS_STATEDIR" description:"The location dr-provision should store running state in"`
	Oneshot    bool   `short:"1" long:"oneshot" env:"RS_ONESHOT" description:"Do not wait for additional tasks"`
	ExitOnFail bool   `short:"x" long:"exit-on-failure" env:"RS_EXIT_ON_FAIL" description:"Exit on failure of a task"`
	SkipPower  bool   `short:"p" long:"skipPower" env:"RS_SKIP_POWER" description:"Skip any power cycle actions"`
	RunOnStart bool   `short:"r" long:"runOnStart" env:"RS_RUN_ON_START" description:"Mark machine as Runnable on startup.  Ignored if Context is set."`
	ConfigFile string `short:"a" long:"agentConfig" env:"RS_AGENT_CONFIG" description:"Configuration for the agent. Overrides settings from the command line and environment" default:"agent.yaml"`
}

var machineAgent *agent.Agent

const (
	serviceName = "drp-agent"
)

func stopAgent() error {
	return machineAgent.Kill()
}

func createAgent(options opts, logger io.Writer) error {
	machine := &models.Machine{}
	session, err := api.TokenSession(options.Endpoint, options.Token)
	if err != nil {
		return fmt.Errorf("Failed connection to %s: %v", options.Endpoint, err)
	}
	if err := session.FillModel(machine, options.MachineID); err != nil {
		return fmt.Errorf("Failed to fetch machine %s: %v", options.MachineID, err)
	}
	if machine.Context == "" && options.Context == "" && options.RunOnStart {
		p := jsonpatch2.Patch{
			{Op: "replace", Path: "/Runnable", Value: true},
		}
		if err := session.Req().Patch(p).UrlForM(machine).Do(machine); err != nil {
			return fmt.Errorf("Failed to mark machine %s runnable: %v\n", options.MachineID, err)
		}
	}
	machineAgent, err = agent.New(session, machine, options.Oneshot, options.ExitOnFail, !options.SkipPower, logger)
	if err != nil {
		return fmt.Errorf("Error starting Agent: %v", err)
	}
	machineAgent.StateLoc(options.StateDir).Context(options.Context)
	return nil
}

func main() {
	var options opts
	var stateLoc string
	switch runtime.GOOS {
	case "windows":
		stateLoc = os.ExpandEnv("${APPDATA}/drp-agent")
	default:
		stateLoc = "/var/lib/drp-agent"
	}

	parser := flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	if options.StateDir == "" {
		options.StateDir = stateLoc
	}
	if options.ConfigFile != "" {
		if !filepath.IsAbs(options.ConfigFile) {
			options.ConfigFile = path.Join(stateLoc, options.ConfigFile)
		}
		if fi, err := os.Stat(options.ConfigFile); err == nil && fi.Mode().IsRegular() {
			buf, err := ioutil.ReadFile(options.ConfigFile)
			if err != nil {
				log.Fatalf("Failed to read config file %s: %v", options.ConfigFile, err)
			}
			if err := models.DecodeYaml(buf, &options); err != nil {
				log.Fatalf("Failed to parse config file %s: %v", options.ConfigFile, err)
			}
		}
	}
	serve(options)
	os.Exit(0)

}

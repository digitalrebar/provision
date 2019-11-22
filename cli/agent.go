package cli

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/VictorLowther/jsonpatch2"

	"github.com/digitalrebar/provision/v4/api"

	"github.com/digitalrebar/provision/v4/models"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

type agentOpts struct {
	Endpoints    string
	Token        string
	MachineID    string
	Context      string
	Oneshot      bool
	ExitOnFail   bool
	SkipPower    bool
	SkipRunnable bool
}

var agentScratchConfig = `---
# Endpoints is an (optional) comma-separated list of dr-provision
# servers that the agent should try to connect to.
# The first one that succeeds will be the one that we use
#
# You must specify at least one endpoint

Endpoints: "https://endpoint1:8092,https://endpoint2:8092"

# Token is the authentication token that the Agent should use.
# It must be a machine-specific token with a lifetime that is 
# significantly longer than the expected time between reprovisions
# of this system.  If you are going to generate this config file
# with a template from a task running during OS install (strongly
# recommended), then {{.GenerateInfiniteToken}} will make a suitable
# token, otherwise the following drpcli command will:
#
# drpcli machines token <machine-uuid> --ttl 3y
#
# That will generate a machine token that expires in 3 years.
#
# You must specify a token

Token: 'base64-encoded token string'

# MachineID is the UUID of the machine.
#
# You must specify a machine UUID

MachineID: 'uuid-of-machine'

# Context is the machine context that this runner should pay attention to.
# When running this agent as the maun agent running tasks for the
# Machine, it should be left blank.

Context: ''

# SkipRunnable tells the agent to not mark the machine as Runnable
# whenever it starts up

SkipRunnable: false
`

type agentProg struct {
	exe          string
	stateLoc     string
	cmd          *exec.Cmd
	opts         agentOpts
	logPipe      io.ReadCloser
	shuttingDown bool
}

func (a *agentProg) Start(s service.Service) error {
	slog, err := s.Logger(nil)
	if err != nil {
		return fmt.Errorf("Error opening logger: %v", err)
	}
	Session, err = sessionOrError(a.opts.Token, strings.Split(a.opts.Endpoints, ","))
	if err != nil {
		return fmt.Errorf("Unable to create session: %v", err)
	}
	proxySock := path.Join(a.stateLoc, "agent.sock")
	os.Remove(proxySock)
	if err := Session.MakeProxy(proxySock); err != nil {
		return fmt.Errorf("Unable to create proxy socket: %v", err)
	}
	if !a.opts.SkipRunnable && a.opts.Context == "" {
		machine := &models.Machine{}
		p := jsonpatch2.Patch{
			{Op: "replace", Path: "/Runnable", Value: true},
		}
		if err := Session.Req().Patch(p).UrlFor("machines", a.opts.MachineID).Do(machine); err != nil {
			return fmt.Errorf("Failed to mark machine %s runnable: %v\n", a.opts.MachineID, err)
		}
	}
	up := &sync.WaitGroup{}
	up.Add(1)
	go func() {
		started := false
		for !a.shuttingDown {
			in, out := io.Pipe()
			a.logPipe = in
			a.cmd.Stderr, a.cmd.Stdout = out, out
			go func(rdr io.Reader) {
				scanner := bufio.NewScanner(rdr)
				for scanner.Scan() {
					slog.Info(scanner.Text())
				}
			}(in)
			err = a.cmd.Start()
			if !started {
				started = true
				up.Done()
				if err != nil {
					in.Close()
					out.Close()
					return
				}
			}
			err := a.cmd.Wait()
			if a.shuttingDown {
				return
			}
			if err != nil {
				slog.Errorf("Agent runner exited prematurely with failure: %v", err)
			} else {
				slog.Errorf("Agent runner exited prematurely!")
			}
			slog.Errorf("Will try again in 5 seconds")
			in.Close()
			out.Close()
			time.Sleep(5 * time.Second)
		}
	}()
	up.Wait()
	return err
}

func (a *agentProg) Stop(s service.Service) error {
	if a == nil {
		log.Panicf("Nil a")
	}
	if a.cmd == nil {
		log.Panicf("Nil cmd")
	}
	if a.cmd.Process == nil {
		log.Panicf("Nil Process")
	}
	if a.cmd.ProcessState == nil || a.cmd.ProcessState.Exited() == false {
		a.shuttingDown = true
		a.cmd.Process.Kill()
	}
	a.logPipe.Close()
	return nil
}

func sessionOrError(token string, endpoints []string) (res *api.Client, err error) {
	for _, endpoint := range endpoints {
		res, err = api.TokenSession(endpoint, token)
		if err == nil {
			return
		}
	}
	return
}

var agentHandler = &cobra.Command{
	Use:   "agent [operation]",
	Short: "Manage drpcli running as an agent",
	Long:  "Use this command to install, remove, stop, start, restart drpcli running as a task runner",
	Args: func(c *cobra.Command, args []string) error {
		if !service.Interactive() {
			return nil
		}
		if len(args) != 1 {
			return fmt.Errorf("%v needs at least 1 argument", c.UseLine())
		}
		switch args[0] {
		case "install", "remove", "stop", "start", "restart", "status":
		default:
			return fmt.Errorf("Unknown agent command %s. Try one of install,remove,stop,start,restart,status", args[0])
		}
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		var stateLoc string
		options := agentOpts{}
		switch runtime.GOOS {
		case "windows":
			stateLoc = `C:/Windows/system32/configs/systemprofile/AppData/Local/rackn/drp-agent`
		default:
			stateLoc = "/var/lib/drp-agent"
		}
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("Unable to determine executable name: %v", err)
		}
		cfgFileName := path.Join(stateLoc, "agent-cfg.yml")
		serviceConfig := &service.Config{
			Name:        "drp-agent",
			DisplayName: "DigitalRebar Provision Agent",
			Description: "DigitalRebar Provision Agent",
			Arguments:   []string{"agent"},
			Executable:  exePath,
		}
		system := service.ChosenSystem()
		switch system.String() {
		case "linux-systemd":
			serviceConfig.Dependencies = []string{
				"Wants=network-online.target network.target",
				"After=network-online.target network.target",
			}
		case "windows-service":
			serviceConfig.Dependencies = []string{
				"Tcpip",
				"Dnscache",
				"LanmanServer",
			}
		}
		prog := &agentProg{
			exe:      exePath,
			stateLoc: stateLoc,
		}
		if !service.Interactive() {
			fi, err := os.Open(cfgFileName)
			if err != nil {
				return fmt.Errorf("Failed to open config file %s: %v", cfgFileName, err)
			}
			defer fi.Close()
			buf, err := ioutil.ReadAll(fi)
			if err != nil {
				return fmt.Errorf("Error reading config file %s: %v", cfgFileName, err)
			}
			if err := models.DecodeYaml(buf, &options); err != nil {
				return fmt.Errorf("Error loading config file %s: %v", cfgFileName, err)
			}
			prog.opts = options
			prog.cmd = exec.Command(exePath, "machines", "processjobs", "--stateDir", stateLoc)
			prog.cmd.Env = append(os.Environ(),
				"RS_ENDPOINTS="+options.Endpoints,
				"RS_TOKEN="+options.Token,
				"RS_UUID="+options.MachineID,
				"RS_CONTEXT="+options.Context)
			svc, err := service.New(prog, serviceConfig)
			if err != nil {
				return fmt.Errorf("Error creating service: %v", err)

			}
			return svc.Run()
		}

		svc, err := service.New(prog, serviceConfig)
		if err != nil {
			return fmt.Errorf("Error creating service: %v", err)

		}
		switch args[0] {
		case "install":
			if err := os.MkdirAll(stateLoc, 0700); err != nil {
				log.Fatalf("Error creating state directory for the dr-provision agent: %v", err)
			}
			fi, err := os.Open(cfgFileName)
			if err != nil {
				fi, err = os.Create(cfgFileName)
				if err != nil {
					log.Fatalf("Unable to create skeleton config file %s: %v", cfgFileName, err)
				}
				defer fi.Close()
				agentEndpoint := os.Getenv("RS_ENDPOINT")
				agentUUID := os.Getenv("RS_UUID")
				agentToken := os.Getenv("RS_TOKEN")
				if agentEndpoint != "" && agentUUID != "" && agentToken != "" {
					log.Printf("RS_* environmment variables present, attempting agent config auto-generation")
					Session, err = api.TokenSession(agentEndpoint, agentToken)
					if err == nil {
						machineToken := &models.UserToken{}
						if err = Session.Req().UrlFor("machines", agentUUID, "token").Params("ttl", "3y").Do(machineToken); err == nil {
							options.Endpoints = agentEndpoint
							options.MachineID = agentUUID
							options.Token = machineToken.Token
							buf, err := api.Pretty("yaml", options)
							if err == nil {
								if _, err = fi.Write(buf); err == nil {
									if err := svc.Install(); err != nil {
										log.Fatalf("Error installing service: %v", err)
									}
									log.Printf("Service installed with auto-generated config %s", cfgFileName)
									return nil
								}
							}
						}
					}
					log.Printf("Unable to auto-fill %s, using scratch config instead", cfgFileName)
				}

				if b, err := fi.Write([]byte(agentScratchConfig)); err != nil || b != len(agentScratchConfig) {
					log.Fatalf("Unable to write sample agent config to %s", cfgFileName)
				}
				log.Fatalf(`Skeleton config file created at %v. 
Please fill it out with final running values, and rerun drpcli agent install`, cfgFileName)
			}
			buf, err := ioutil.ReadAll(fi)
			fi.Close()
			if err != nil {
				log.Fatalf("Error reading config file %s: %v", cfgFileName, err)
			}
			if err := models.DecodeYaml(buf, &options); err != nil {
				log.Fatalf("Error loading config file %s: %v", cfgFileName, err)
			}
			if options.Endpoints == "" || options.MachineID == "" || options.Token == "" {
				log.Fatalf("Config file %s missing a required parameter", cfgFileName)
			}
			Session, err = sessionOrError(options.Token, strings.Split(options.Endpoints, ","))
			if err != nil {
				log.Fatalf("Unable to establish connection to an endpoint listed in %s: %v", cfgFileName, err)
			}
			if exists, err := Session.ExistsModel("machines", options.MachineID); err != nil || !exists {
				log.Fatalf("Failed to verify that machine %s exists on %s", options.MachineID, Session.Endpoint())
			}
			if err := svc.Install(); err != nil {
				log.Fatalf("Error installing service: %v", err)
			}
		case "remove":
			if err := svc.Uninstall(); err != nil {
				log.Fatalf("Error removing service: %v", err)
			}
		case "start":
			if err := svc.Start(); err != nil {
				log.Fatalf("Error starting service: %v", err)
			}
		case "stop":
			if err := svc.Stop(); err != nil {
				log.Fatalf("Error stopping service: %v", err)
			}
		case "restart":
			if err := svc.Restart(); err != nil {
				log.Fatalf("Error restarting service: %v", err)
			}
		case "status":
			status, err := svc.Status()
			if err != nil {
				log.Fatalf("Error getting service status: %v", err)
			}
			switch status {
			case service.StatusRunning:
				prettyPrint("running")
			case service.StatusStopped:
				prettyPrint("stopped")
			default:
				log.Fatalf("Unknown status %v", status)
			}
		}
		return nil
	},
}

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

	"github.com/digitalrebar/service"
	"github.com/spf13/cobra"
)

type agentOpts struct {
	Endpoints       string
	Token           string
	MachineID       string
	Context         string
	Oneshot         bool
	ExitOnFail      bool
	SkipPower       bool
	SkipRunnable    bool
	AllowAutoUpdate bool
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
# The token itself is the base64-encoded string in the Token field
# of the returned JSON object.
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

# AllowAutoUpdate allows the agent to try and automatically update itself
# as needed whenever it is starting up.  This does not work on Windows.

AllowAutoUpdate: true
`

type agentProg struct {
	exe          string
	stateLoc     string
	serviceType  string
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
			switch a.serviceType {
			case "linux-systemd", "linux-upstart":
				a.cmd.Stderr = os.Stderr
				a.cmd.Stdout = os.Stdout
			default:
				in, out, err := os.Pipe()
				if err != nil {
					slog.Errorf("Error setting up logging: %v", err)
					slog.Errorf("Will try again in 5 seconds")
					time.Sleep(5 * time.Second)
					continue
				}
				a.logPipe = in
				a.cmd.Stderr, a.cmd.Stdout = out, out
				go func(rdr, wrt io.ReadCloser) {
					defer wrt.Close()
					scanner := bufio.NewScanner(rdr)
					for scanner.Scan() {
						slog.Info(scanner.Text())
					}
				}(in, out)
			}
			err = a.cmd.Start()
			if !started {
				started = true
				up.Done()
				if err != nil {
					if a.logPipe != nil {
						a.logPipe.Close()
					}
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
			if a.logPipe != nil {
				a.logPipe.Close()
			}
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
	if a.logPipe != nil {
		a.logPipe.Close()
	}
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
		if len(args) > 1 {
			return fmt.Errorf("%v needs at at most 1 argument", c.UseLine())
		}
		if len(args) == 0 {
			return nil
		}
		switch args[0] {
		case "install", "remove", "stop", "start", "restart", "status":
		default:
			return fmt.Errorf("Unknown agent command %s. Try one of install,remove,stop,start,restart,status", args[0])
		}
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		stateLoc := DefaultStateLoc
		options := agentOpts{}
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
			exe:         exePath,
			stateLoc:    stateLoc,
			serviceType: system.String(),
		}
		if len(args) == 0 {
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
			// Auto-update can only happen on non-Windows, because we need to replace the running binary.
			func() {
				if runtime.GOOS == "windows" {
					log.Printf("Unable to update on Windows")
					return
				}
				if !options.AllowAutoUpdate {
					log.Printf("Auto update disabled by config directive")
					return
				}
				session, err := sessionOrError(options.Token, strings.Split(options.Endpoints, ","))
				if err != nil {
					log.Printf("No session")
					return
				}
				info, err := session.Info()
				if err != nil || !info.HasFeature("agent-auto-update") {
					log.Printf("Missing agent-auto-update from dr-provision, will not auto update")
					return
				}
				remoteName := fmt.Sprintf("drpcli.%s.%s", runtime.GOARCH, runtime.GOOS)
				sum, err := session.GetBlobSum("files", remoteName)
				if err != nil {
					log.Printf("No blob sum for %s", remoteName)
					return
				}
				exe, err := os.Open(exePath)
				if err != nil {
					log.Printf("No exe at %s", exePath)
					return
				}
				mta := &models.ModTimeSha{}
				if _, err := mta.Regenerate(exe); err == nil && mta.String() == sum {
					log.Printf("Sums for %s already match. No update needed", exePath)
					exe.Close()
					return
				}
				exe.Close()
				tmpName := exePath + ".new"
				exe, err = os.OpenFile(tmpName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
				if err != nil {
					log.Printf("Error opening tmp %s: %v", tmpName, err)
					return
				}
				defer os.Remove(tmpName)
				if err := session.GetBlob(exe, "files", remoteName); err != nil {
					log.Printf("Error getting %s: %v", remoteName, err)
					exe.Close()
					return
				}
				session.Close()
				exe.Close()
				testCmd := exec.Command(tmpName, "version")
				if err := testCmd.Run(); err != nil {
					log.Printf("Error validating %s version: %v", tmpName, err)
					return
				}
				if err := os.Rename(tmpName, exePath); err != nil {
					log.Printf("Error renaming %s to %s: %v", tmpName, exePath, err)
				} else {
					log.Printf("%s updated from %s to %s", exePath, mta.String(), sum)
				}
			}()
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
					log.Printf("Unable to auto-fill %s, using scratch config instead: %v", cfgFileName, err)
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

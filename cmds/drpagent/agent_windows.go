// +build windows

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sync/semaphore"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func mgrAndService(name string) (*mgr.Mgr, *mgr.Service, error) {
	m, err := mgr.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the service manager: %v", err)
	}
	s, err := m.OpenService(name)
	if err != nil {
		m.Disconnect()
	}
	return m, s, err
}

func installMyself(name string) error {
	exeName, err := filepath.Abs(os.Args[0])
	if err != nil {
		return fmt.Errorf("Failed to determine path to myself: %v", err)
	}
	m, s, err := mgrAndService(name)
	defer m.Disconnect()
	if err == nil {
		s.Close()
		return fmt.Errorf("%s already registered as a service", name)
	}
	s, err = m.CreateService(
		name,
		exeName,
		mgr.Config{
			DisplayName:      "DigitalRebar Provision Agent",
			StartType:        mgr.StartAutomatic,
			DelayedAutoStart: true,
			ErrorControl:     mgr.ErrorNormal,
		},
		"is",
		"auto-started")
	if err != nil {
		return fmt.Errorf("Error creating service %s", name)
	}
	defer s.Close()
	if err := eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info); err != nil {
		s.Delete()
		return fmt.Errorf("Failed to create event log for %s: %v", name, err)
	}
	return nil
}

func removeMyself(name string) error {
	m, s, err := mgrAndService(name)
	defer m.Disconnect()
	if err != nil {
		return fmt.Errorf("Service %s not installed: %v", name, err)
	}
	defer s.Close()
	if err := s.Delete(); err != nil {
		return fmt.Errorf("Error deleting service %s: %v", name, err)
	}
	if err := eventlog.Remove(name); err != nil {
		return fmt.Errorf("Error removing event log for service %s: %v", name, err)
	}
	return nil
}

func startMyself(name string) error {
	m, s, err := mgrAndService(name)
	defer m.Disconnect()
	if err != nil {
		return fmt.Errorf("Service %s not installed: %v", name, err)
	}
	defer s.Close()
	if err := s.Start("is", "manual-start"); err != nil {
		return fmt.Errorf("Failed to start service %s: %v", name, err)
	}
	return nil
}

func stopMyself(name string) error {
	m, s, err := mgrAndService(name)
	defer m.Disconnect()
	if err != nil {
		return fmt.Errorf("Service %s not installed: %v", name, err)
	}
	defer s.Close()
	status, err := s.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("Could not send stop to %s: %v", name, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != svc.Stopped {
		if time.Now().After(timeout) {
			return fmt.Errorf("Service failed to stop after waiting 10 seconds")
		}
		time.Sleep(time.Second)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("Could not fetch status of %s: %v", name, err)
		}
	}
	return nil
}

type myLog struct {
	*eventlog.Log
	options opts
	name    string
}

func (m *myLog) Write(buf []byte) (int, error) {
	err := m.Info(1, string(buf))
	return len(buf), err
}

func (m *myLog) Execute(args []string,
	changeReqs <-chan svc.ChangeRequest,
	changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	var runErr error
	m.Info(1, fmt.Sprintf("Starting %s", m.name))
	err := createAgent(m.options, m)
	if err != nil {
		m.Error(1, fmt.Sprintf("%v", err))
		return
	}
	sema := semaphore.NewWeighted(1)
	sema.Acquire(context.Background(), 1)
	go func() {
		runErr = machineAgent.Run()
		sema.Release(1)
	}()
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
loop:
	for {
		change := <-changeReqs
		switch change.Cmd {
		case svc.Interrogate:
			newStatus := change.CurrentStatus
			if sema.TryAcquire(1) {
				if runErr != nil {
					m.Error(1, fmt.Sprintf("Unexpected agent shutdown: %v", runErr))
					m.Error(1, fmt.Sprintf("Service will exit and should restart"))
					errno = 1
					return
				}
				sema.Release(1)
				break loop
			}
			changes <- newStatus
			time.Sleep(100 * time.Millisecond)
			changes <- newStatus
		case svc.Stop, svc.Shutdown:
			stopAgent()
			break loop
		default:
			m.Error(1, fmt.Sprintf("Unexpected control request %v", change))

		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runMyself(name string, options opts) {
	elog, err := eventlog.Open(name)
	if err != nil {
		return
	}
	defer elog.Close()
	l := &myLog{elog, options, name}
	svc.Run(name, l)
}

func serve(options opts) {
	svcName := "drp-agent"
	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("Cannot determine if we are running interactively or not")
	}
	if isInteractive {
		if len(os.Args) == 1 {
			fmt.Print(`Valid commands are:
install  Installs this command as the drp-agent service.
remove   Removes the drp-agent service.
start    Starts the drp-agent service.  Must be installed first.
stop     Stops the drp-agent service.   Must be started first.
`)
			os.Exit(1)
		}
		switch os.Args[1] {
		case "install":
			err = installMyself(svcName)
		case "remove":
			err = removeMyself(svcName)
		case "start":
			err = startMyself(svcName)
		case "stop":
			err = stopMyself(svcName)
		default:
			log.Fatalf("Unknown command %s", os.Args[1])
		}
		if err != nil {
			log.Fatalf("Error running %s: %v", os.Args[1], err)
		}
		return
	}
	runMyself(svcName, options)
}

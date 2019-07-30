package test

import (
	"os"
	"os/exec"
	"path"
	"syscall"
)

var (
	server *exec.Cmd
)

func StartServer(tmpDir string) error {
	os.Setenv("RS_TOKEN_PATH", path.Join(tmpDir, "tokens"))
	os.Setenv("RS_ENDPOINT", "https://127.0.0.1:10001")

	server = exec.Command("dr-provision",
		"--base-root", tmpDir,
		"--tls-key", tmpDir+"/server.key",
		"--tls-cert", tmpDir+"/server.crt",
		"--api-port", "10001",
		"--static-port", "10002",
		"--tftp-port", "10003",
		"--dhcp-port", "10004",
		"--binl-port", "10005",
		"--metrics-port", "10006",
		"--fake-pinger",
		"--drp-id", "Fred",
		"--plugin-comm-root", tmpDir,
		"--backend", "memory:///")
	server.Stdout = os.Stderr
	server.Stderr = os.Stderr
	server.Dir = tmpDir
	server.Env = os.Environ()
	server.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGTERM}
	return server.Start()
}

func StopServer() error {
	server.Process.Signal(os.Kill)
	return server.Wait()
}

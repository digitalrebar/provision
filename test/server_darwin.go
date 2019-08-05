package test

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

var (
	server *exec.Cmd
)

func StartServer(tmpDir string, basePort int) error {
	apiPort := fmt.Sprintf("%d", basePort)
	staticPort := fmt.Sprintf("%d", basePort+1)
	tftpPort := fmt.Sprintf("%d", basePort+2)
	dhcpPort := fmt.Sprintf("%d", basePort+3)
	binlPort := fmt.Sprintf("%d", basePort+4)
	metricPort := fmt.Sprintf("%d", basePort+5)

	os.Setenv("RS_TOKEN_PATH", path.Join(tmpDir, "tokens"))
	os.Setenv("RS_ENDPOINT", fmt.Sprintf("https://127.0.0.1:%s", apiPort))

	server = exec.Command("dr-provision",
		"--base-root", tmpDir,
		"--tls-key", tmpDir+"/server.key",
		"--tls-cert", tmpDir+"/server.crt",
		"--api-port", apiPort,
		"--static-port", staticPort,
		"--tftp-port", tftpPort,
		"--dhcp-port", dhcpPort,
		"--binl-port", binlPort,
		"--metrics-port", metricPort,
		"--fake-pinger",
		"--drp-id", "Fred",
		"--plugin-comm-root", tmpDir,
		"--backend", "memory:///")
	server.Stdout = os.Stderr
	server.Stderr = os.Stderr
	server.Dir = tmpDir
	server.Env = os.Environ()
	return server.Start()
}

func StopServer() error {
	server.Process.Signal(os.Kill)
	server.Process.Kill()
	return nil
}

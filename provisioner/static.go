package provisioner

import (
	"net"
	"net/http"
)

func ServeStatic(listenAt, fsPath string) error {
	conn, err := net.Listen("tcp", listenAt)
	if err != nil {
		return err
	}
	fs := http.FileServer(http.Dir(fsPath))
	http.Handle("/", fs)
	return http.Serve(conn, nil)
}

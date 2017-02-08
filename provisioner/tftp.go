package provisioner

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/pin/tftp"
)

func handleTftpRead(filename string, rf io.ReaderFrom) error {
	p := filepath.Join(ProvOpts.FileRoot, filename)
	p = filepath.Clean(p)
	if !strings.HasPrefix(p, ProvOpts.FileRoot) {
		err := fmt.Errorf("Filename %s tries to escape root %s", filename, ProvOpts.FileRoot)
		log.Println(err)
		return err
	}
	log.Printf("Sending %s from %s", filename, p)
	file, err := os.Open(p)
	if err != nil {
		log.Println(err)
		return err
	}
	if t, ok := rf.(tftp.OutgoingTransfer); ok {
		if fi, err := file.Stat(); err == nil {
			t.SetSize(fi.Size())
		}
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("%d bytes sent\n", n)
	return nil
}

func ServeTftp(listen string) error {
	a, err := net.ResolveUDPAddr("udp", listen)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", a)
	if err != nil {
		return err
	}
	svr := tftp.NewServer(handleTftpRead, nil)
	go svr.Serve(conn)
	return nil
}

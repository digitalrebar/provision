// +build !windows,!plan9

package api

import "os"
import "golang.org/x/sys/unix"

func mktd(p string) error {
	umask := unix.Umask(0)
	defer unix.Umask(umask)
	return os.MkdirAll(p, 0777|os.ModeSticky)
}

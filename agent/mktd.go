// +build !windows,!plan9

package agent

import (
	"os"
	"sync"
)
import "golang.org/x/sys/unix"

var mtdLock = &sync.Mutex{}

func mktd(p string) error {
	mtdLock.Lock()
	defer mtdLock.Unlock()
	umask := unix.Umask(0)
	defer unix.Umask(umask)
	return os.MkdirAll(p, 0777|os.ModeSticky)
}

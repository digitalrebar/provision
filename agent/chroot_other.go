// +build !linux

package agent

import (
	"fmt"
	"os/exec"
	"runtime"
)

func (r *runner) enterChroot(cmd *exec.Cmd) error {
	if r.chrootDir != "" {
		return fmt.Errorf("enterChroot not supported on %v", runtime.GOOS)
	}
	return nil
}

func (r *runner) exitChroot() {}

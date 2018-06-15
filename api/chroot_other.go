// +build !linux

package api

import (
	"fmt"
	"os/exec"
	"runtime"
)

func (r *TaskRunner) enterChroot(cmd *exec.Cmd) error {
	if r.chrootDir != "" {
		return fmt.Errorf("enterChroot not supported on %v", runtime.GOOS)
	}
	return nil
}

func (r *TaskRunner) exitChroot() {}

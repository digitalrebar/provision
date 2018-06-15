// +build linux

package api

import (
	"os"
	"os/exec"
	"path"
	"syscall"
)

func bindMount(newRoots string, srcFS ...string) error {
	if len(srcFS) == 0 {
		return nil
	}
	tgt := path.Join(newRoots, srcFS[0])
	if err := os.MkdirAll(tgt, os.ModePerm); err != nil {
		return err
	}
	if err := syscall.Mount(srcFS[0], tgt, "", syscall.MS_BIND, ""); err != nil {
		return err
	}
	if err := bindMount(newRoots, srcFS[1:]...); err != nil {
		syscall.Unmount(tgt, 0)
		return err
	}
	return nil
}

func (r *TaskRunner) bindFSes() []string {
	return []string{"/proc", "/sys", "/dev", "/dev/pts", r.agentDir}
}

func (r *TaskRunner) enterChroot(cmd *exec.Cmd) error {
	if r.chrootDir == "" {
		return nil
	}
	if err := bindMount(r.chrootDir, r.bindFSes()...); err != nil {
		return err
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Chroot: r.chrootDir}
	return nil
}

func (r *TaskRunner) exitChroot() {
	if r.chrootDir == "" {
		return
	}
	fses := r.bindFSes()
	for i := len(fses) - 1; i > -1; i-- {
		syscall.Unmount(path.Join(r.chrootDir, fses[i]), 0)
	}
}

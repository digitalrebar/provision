// +build windows plan9

package agent

import "os"

func mktd(p string) error {
	return os.MkdirAll(p, 01777)
}

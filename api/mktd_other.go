// +build windows plan9

package api

import "os"

func mktd(p string) error {
	return os.MkdirAll(p, 01777)
}

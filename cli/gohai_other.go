// +build !linux

package cli

// for darwin help one day: https://github.com/cavaliercoder/dmidecode-osx/blob/master/dmidecode.c
// for windows, similar dmidecode code.

func gohai() error {
	infos := map[string]string{"Error": "Unsupported platform"}
	prettyPrint(infos)
	return nil
}

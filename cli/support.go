/*
Copyright © 2020 RackN <support@rackn.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cli

import (
"archive/zip"
"bytes"
"fmt"
"github.com/spf13/cobra"
"io"
"io/ioutil"
"log"
"os"
"os/exec"
"runtime"
"strings"
"time"
)

// bundleCmd represents the bundle command
var drUser string
var drBase string
var since string
var extraDirs string

func bundleCmds()  *cobra.Command {
	cmd := &cobra.Command{
		Use:   "support",
		Short: "Access commands related to RackN Tech Support",
	}

	bundleCmd := &cobra.Command{
		Use:   "bundle",
		Short: "Create a support bundle for the RackN engineering team.",
		Long: `Create a support bundle for the RackN engineering team.
	This command is currently only supported on a Linux host and
	expects to be running on the drp endpoint.

	By default the command will run:
		journalctl -u dr-provision --since yesterday
	It captures that output and puts it into a file.
	Next we take the contents of /var/lib/dr-provision
	excluding some folders and add them along with the
	log output to a zip file that can be sent to support@rackn.com

	If your drp endpoint runs as some other user you can set the user with the --dr-user flag.
	If your drp endpoint has a different base dir than /var/lib/dr-provision
	you can set that with the --drp-basedir flag.
	If you need to include additional directories that are excluded by default above
	you can add them with  --extra-dirs This is only needed if directed by support.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tmpBuf := new(bytes.Buffer)
		w := zip.NewWriter(tmpBuf)
		var command = "journalctl -u " + drUser + " --since " + since
		outBuf, errBuf, err := runCommnd(command)
		if err != nil {
			fmt.Printf("An error happened trying to run %s. Got: %s", command, errBuf.String() )
			return err
		}
		f, err := w.Create("rackn_bundle/journal-output.txt")
		if err != nil {
			return err
		}
		f.Write(outBuf.Bytes())
		// Loop through the directories we want to add from the base
		if ! strings.HasSuffix(drBase, "/") {
			drBase = drBase + "/"
		}
		dirsToGet := []string {"wal", "digitalrebar", "secrets"}
		if len(extraDirs) > 0 {
			extras := strings.Split(extraDirs, ",")
			if extras != nil && len(extras) > 0 {
				dirsToGet = append(dirsToGet, extras...)
			}
		}
		for _, d := range dirsToGet {
			err = AddDirToZip(w,drBase + d + "/")
			if err != nil {
				return err
			}
		}

		if w.Close() != nil {
			return err
		}

		t := time.Now()
		fname := fmt.Sprintf("drp-support-bundle-%d-%02d-%02d-%02d-%02d-%02d.zip",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
		err = ioutil.WriteFile(fname, tmpBuf.Bytes(), 0644)
		if err != nil {
			return err
		}
		return nil
		},
	}
	bundleCmd.Flags().StringVarP(&extraDirs, "extra-dirs", "", "", "extra-dirs job-logs,saas-content,ux,plugins")
	bundleCmd.Flags().StringVarP(&since, "since", "", "yesterday", "since 'something valid that journalctl supports'")
	bundleCmd.Flags().StringVarP(&drUser, "dr-user", "", "dr-provision", "dr-user dr-provision")
	bundleCmd.Flags().StringVarP(&drBase, "drp-basedir", "", "/var/lib/dr-provision/", "drp-basedir /var/lib/dr-provision")
	cmd.AddCommand(bundleCmd)

	return cmd
}



func init() {
	addRegistrar(func(c *cobra.Command) { c.AddCommand(bundleCmds()) })
}

// run a system command and return the std output as a bytes buffer
func runCommnd(cmdStr string) (outBuffer bytes.Buffer, errBuf bytes.Buffer, err error) {
	if runtime.GOOS != "linux" {
		// Inform the user their platform is currently unsupported.
		fmt.Println("Currently Linux is the only supported platform for this feature.")
		os.Exit(1)
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	c := strings.Split(cmdStr, " ")
	if err := checkToolExists(c[0]); err != nil {
		// fail
		return stdoutBuf, stderrBuf, err
	}
	cmd := exec.Command(c[0], c[1:]...)
	cmd.Stdout = io.MultiWriter(&stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return  stdoutBuf, stderrBuf, err
	}
	return stdoutBuf, errBuf, err
}

// verify the command bring run is on the system in the path
func checkToolExists(t string) error {
	_,err := exec.LookPath(t)
	if err != nil {
		fmt.Printf("didn't find %s in path\n", t)
		return err
	}
	return nil
}

// given a top level dir this will walk it adding files and folders.
// to the archive
func AddDirToZip(zipWriter *zip.Writer, dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		if strings.HasSuffix(dirname, "secrets/") { return nil }
		if strings.HasSuffix(dirname, "digitalrebar/") { return nil }
		fmt.Println("encountered an error ", err)
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			err = AddFileToZip(zipWriter, dirname + file.Name())
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// Recurse
			newBase := dirname + file.Name() + "/"
			AddDirToZip(zipWriter, newBase)
		}
	}
	return err
}

// adds an individual file to the archive
func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	fileToZip.Close()
	return err
}

package cli

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

var installSkipDownloadIsos = true

func init() {
	addRegistrar(registerBootEnv)
}

func registerBootEnv(app *cobra.Command) {
	op := &ops{
		name:       "bootenvs",
		singleName: "bootenv",
		example:    func() models.Model { return &models.BootEnv{} },
	}
	installCmd := &cobra.Command{
		Use:   "install [bootenvFile] [isoPath]",
		Short: "Install a bootenv along with everything it requires",
		Long: `bootenvs install assumes a directory with two subdirectories:
   bootenvs/
   templates/

bootenvs must contain [bootenvFile]
templates must contain any templates that the requested bootenv refers to.

bootenvs install will try to upload any required ISOs if they are not already
present in DigitalRebar Provision.  If [isoPath] is specified, it will use that
directory to to check and download ISOs into, otherwise it will use isos/  If the
ISO is not present, we will try to download it if the bootenv specifies a location
to download the ISO from.  If we cannot find an ISO to upload, then the bootenv
will still be uploaded, but it will not be available until the ISO is uploaded
using isos upload.git `,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%v needs at least 1 arg", c.UseLine())
			}
			if len(args) > 2 {
				return fmt.Errorf("%v has Too many args", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			isoCache := "isos"
			if len(args) == 2 {
				isoCache = args[1]
			}
			bootEnv, err := session.InstallBootEnvFromFile(args[0])
			if err != nil {
				return generateError(err, "Failed to install bootenv")
			}
			if bootEnv.OS.IsoFile != "" && !installSkipDownloadIsos {
				if err = os.MkdirAll(isoCache, 0755); err != nil {
					return fmt.Errorf("Error ensuring ISO cache exists: %s", err)
				}
				isoPath := path.Join(isoCache, bootEnv.OS.IsoFile)
				if err := session.InstallISOForBootenv(bootEnv, isoPath, !installSkipDownloadIsos); err != nil {
					return generateError(err, "Error uploading %s", isoPath)
				}
			}
			return prettyPrint(bootEnv)
		},
	}
	installCmd.Flags().BoolVar(&installSkipDownloadIsos,
		"skip-download",
		false,
		"Whether to try to download ISOs from their upstream")
	op.addCommand(installCmd)
	op.addCommand(&cobra.Command{
		Use:   "uploadiso [id]",
		Short: "This will attempt to upload the ISO from the specified ISO URL.",
		Long: `This will attempt to upload the ISO from the specified ISO URL.
It will attempt to perform a direct copy without saving the ISO locally.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			bootEnv := &models.BootEnv{}
			if err := session.FillModel(bootEnv, args[0]); err != nil {
				return generateError(err, "Failed to fetch %v: %v", op.singleName, args[0])
			}
			if bootEnv.OS.IsoFile == "" {
				return fmt.Errorf("BootEnv %s does not require an iso", bootEnv.Name)
			}
			if bootEnv.OS.IsoUrl == "" {
				return fmt.Errorf("Unable to automatically download iso for %s", bootEnv.Name)
			}
			isoDlResp, err := http.Get(bootEnv.OS.IsoUrl)
			if err != nil {
				return fmt.Errorf("Unable to connect to %s: %v", bootEnv.OS.IsoUrl, err)
			}
			defer isoDlResp.Body.Close()
			if isoDlResp.StatusCode >= 300 {
				return fmt.Errorf("Unable to initiate download of %s: %s", bootEnv.OS.IsoUrl, isoDlResp.Status)
			}
			if info, err := session.PostBlob(isoDlResp.Body, "isos", bootEnv.OS.IsoFile); err != nil {
				return generateError(err, "Error uploading %s", bootEnv.OS.IsoFile)
			} else {
				return prettyPrint(info)
			}
		},
	})
	op.addCommand(&cobra.Command{
		Use:   "fromAppleNBI [path]",
		Short: "This will attempt to translate an Apple .nbi directory into a bootenv and an archive.",
		Long: `This command translates an Apple .nbi directory into a bootenv .yaml file
that contains apropriate metadata to be handled by the dr-provision NBSP DHCP
handler, and a .tar.gz file that contains the contents of the .nbi directory.

The .nbi directory must have been produced by the Apple System Image Utility
or equivalent tooling, and must contain a valid NBImageInfo.plist file.
The .yaml file containig the bootenv will be named <os>-<version>.yaml,
and the .tar.gz file will contain the contents of the .nbi directory.

Both created files will be left in the current working directory.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			sb, err := os.Stat(path.Join(args[0], "NBImageInfo.plist"))
			if err != nil {
				return fmt.Errorf("Cannot find NBImageInfo.plist in %s: %v", args[0], err)
			}
			if !sb.Mode().IsRegular() {
				return fmt.Errorf("%s is not a normal file", path.Join(args[0], "NBImageInfo.plist"))
			}
			return nil
		},
		RunE: genEnvAndArchiveFromAppleNBI,
	})
	op.command(app)
}

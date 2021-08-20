package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerPluginProvider)
}

func registerPluginProvider(app *cobra.Command) {
	op := &ops{
		name:       "plugin_providers",
		singleName: "plugin_provider",
		example:    func() models.Model { return &models.PluginProvider{} },
		noCreate:   true,
		noUpdate:   true,
	}
	replaceWritable := false
	upload := &cobra.Command{
		Use:   "upload [name] (from [file])",
		Short: "Upload a program to act as a plugin_provider",
		Long: `Uploads a program to act as a plugin_provider.
If the final name of the plugin_provider is the same as the name of the file being uploaded,
then the (from [file]) part may be omitted, and [name] should be the path to the plugin_provider.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 && len(args) != 1 {
				return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			var name, filePath string
			if len(args) == 1 {
				filePath = args[0]
				name = path.Base(args[0])
			} else {
				name = args[0]
				filePath = args[2]
			}
			fi, err := urlOrFileAsReadCloser(filePath)
			if err != nil {
				return fmt.Errorf("Error opening %s: %v", filePath, err)
			}
			defer fi.Close()
			res := &models.PluginProviderUploadInfo{}
			req := Session.Req().Post(fi).UrlFor(op.name, name)
			if replaceWritable {
				req = req.Params("replaceWritable", "true")
			}
			if err := req.Do(res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	}
	upload.Flags().BoolVar(&replaceWritable, "replace-writable", false, "Replace identically named writable objects")
	// Flag deprecated due to standardizing on all hyphenated form for persistent flags.
	upload.Flags().BoolVar(&replaceWritable, "replaceWritable", false, "Replace identically named writable objects")
	upload.Flags().MarkHidden("replaceWritable")
	upload.Flags().MarkDeprecated("replaceWritable", "please use --replace-writable")
	op.addCommand(upload)
	op.addCommand(&cobra.Command{
		Use:     "download [item] to [dest]",
		Aliases: []string{"show", "get"},
		Short:   "Download the plugin_provider named [item] to [dest]",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 || len(args) == 3 {
				return nil
			}
			return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			dest := os.Stdout
			if len(args) == 3 && args[2] != "-" {
				var err error
				dest, err = os.OpenFile(args[2], os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					return fmt.Errorf("Error opening dest file %s: %v", args[2], err)
				}
				defer dest.Close()
			}
			if err := Session.GetBlob(dest, "plugin_providers", args[0], "binary"); err != nil {
				return generateError(err, "Failed to fetch %v: %v", "plugin_providers", args[0])
			}
			return nil
		},
	})
	op.command(app)
}

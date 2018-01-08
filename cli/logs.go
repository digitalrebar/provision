package cli

import "github.com/spf13/cobra"

func registerLog(app *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Access commands relating to logs",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "watch",
		Short: "Watch log entrys as theyt come in real time",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			stream, err := session.Events()
			if err != nil {
				return err
			}
			handle, es, err := stream.Register("log.*.*")
			if err != nil {
				return err
			}
			defer stream.Deregister(handle)
			for {
				evt := <-es
				if evt.Err != nil {
					return err
				}
				prettyPrint(evt.E.Object)
			}
		},
	})
	app.AddCommand(cmd)
}

func init() {
	addRegistrar(registerLog)
}

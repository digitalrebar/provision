package cli

import "github.com/spf13/cobra"

func registerLog(app *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Access commands relating to logs",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "watch",
		Short: "Watch log entrys as they come in real time",
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
	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get the currently buffered log entries from dr-provision",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			res, err := session.Logs()
			if err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	app.AddCommand(cmd)
}

func init() {
	addRegistrar(registerLog)
}

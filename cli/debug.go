package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerDebug)
}

func registerDebug(app *cobra.Command) {
	seconds := 0
	prefix := ""
	debug := &cobra.Command{
		Use:   "debug [type] [target]",
		Short: "Gather [type] of debug information and save it to [target]",
		Long: `This command gathers various different types of runtime profile data from a running
dr-provision server, provided it has the /api/v3/debug or /api/v3/drp_debug.  The types of data that can be gathered are:

    profile: CPU utilization profile information.  Tracks how much CPU time is being used in which
             functions, based on sampling which functions are running every 10 ms.  If the
             --seconds flag is unspecified, profile will gather 30 seconds worth of data.

    trace: Execution trace information, including information on where execution is blocked on
           various types of IO and synchronization primitives.  If the --seconds flag is unspecified,
           trace will gather 1 second of data.

    heap: Memory tracing information for all live data in memory. heap is always point-in-time data.

    heapdump: All live objects in the system.  Use only as directed by support.

    allocs: Memory tracing of all memory that has been allocated since the start of the program
            This includes memory that has been garbage-collected.  alloc is always point-in-time data.

    block: Stack traces of all goroutines that have blocked on synchronization primitives.
           block is always point-in-time data.

    mutex: Stack traces of all holders of contended mutexes.  mutex is always point-in-time data.

    threadcreate: Stack traces of all goroutines that led to the creation of a new OS thread.
                  threadcreate is always point-in-time data.

    goroutine: Stack traces of all current goroutines. goroutine is always point-in-time data.

    index: Returns the indexes of the stacks with the flags of the object.
`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			useseconds := false
			switch args[0] {
			case "profile":
				useseconds = true
				if seconds == 0 {
					seconds = 30
				}
			case "trace":
				useseconds = true
				if seconds == 0 {
					seconds = 1
				}
			case "heap", "allocs", "block", "mutex", "threadcreate", "goroutine", "index", "heapdump":
			default:
				return fmt.Errorf("Unknown debug type %s", args[0])
			}
			target, err := os.Create(args[1])
			if err != nil {
				return err
			}
			defer target.Close()
			req := Session.Req()
			if args[0] == "index" {
				req = req.UrlFor("drp_debug", args[0], prefix)
			} else {
				req = req.UrlFor("debug", args[0])
			}
			if useseconds && seconds > 0 {
				req = req.Params("seconds", fmt.Sprintf("%d", seconds))
			}
			if useseconds {
				log.Printf("Gathering %s debug data for %d seconds to %s", args[0], seconds, args[1])
			} else {
				log.Printf("Gathering %s debug data to %s", args[0], args[1])
			}
			if err := req.Do(target); err != nil {
				return err
			}
			log.Printf("Done")
			return nil
		},
	}
	debug.Flags().IntVar(&seconds, "seconds", 0, "How much debug data to gather, for types that gather data over time.")
	debug.Flags().StringVar(&prefix, "prefix", "", "Limits the index call to just this prefix type.")
	app.AddCommand(debug)
}

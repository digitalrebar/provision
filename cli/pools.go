package cli

import (
	"fmt"
	"strings"

	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerPool)
}

func registerPool(app *cobra.Command) {
	op := &ops{
		name:       "pools",
		singleName: "pool",
		example:    func() models.Model { return &models.Pool{} },
	}

	statusCmd := &cobra.Command{
		Use:   "status [pool]",
		Short: "Get Pool status",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Must provide a pool")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			res := &models.PoolResults{}
			if err := Session.Req().UrlFor("pools", id, "status").Do(&res); err != nil {
				return fmt.Errorf("Failed to get pool status for %s: %v", id, err)
			}
			return prettyPrint(res)
		},
	}
	op.addCommand(statusCmd)

	activeCmd := &cobra.Command{
		Use:   "active",
		Short: "List active pools",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("Does not take parameters")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			res := []string{}
			if err := Session.Req().UrlFor("pools-active").Do(&res); err != nil {
				return fmt.Errorf("Failed to list active pools: %v", err)
			}
			return prettyPrint(res)
		},
	}
	op.addCommand(activeCmd)

	itemCmd := &cobra.Command{
		Use:   "manage",
		Short: "Manage machines in pools",
	}
	// Status
	var allMachines bool
	var minimum, count int
	var addProfiles, removeProfiles, removeParameters string
	var addParameters, newWorkflow, machineList, waitTimeout string
	var sourcePool string

	itemCmd.PersistentFlags().BoolVar(&allMachines, "all-machines", false, "Selects all available machines")
	itemCmd.PersistentFlags().IntVar(&minimum, "minimum", 0, "Minimum number of machines to return - defaults to count")
	itemCmd.PersistentFlags().IntVar(&count, "count", 0, "Count of machines to allocate")
	itemCmd.PersistentFlags().StringVar(&machineList, "machine-list", "", "Comma separated list of machines UUID or Field:Value")
	itemCmd.PersistentFlags().StringVar(&waitTimeout, "wait-timeout", "", "An amount of time to wait for completion in seconds or time string (e.g. 30m)")

	itemCmd.PersistentFlags().StringVar(&newWorkflow, "new-workflow", "", "A workflow to set on the machines")
	itemCmd.PersistentFlags().StringVar(&addProfiles, "add-profiles", "", "Comma separated list of profiles to add to the machine")
	itemCmd.PersistentFlags().StringVar(&removeProfiles, "remove-profiles", "", "Comma separated list of profiles to remove from the machine")
	itemCmd.PersistentFlags().StringVar(&removeParameters, "remove-parameters", "", "Comma separated list of parameters to remove from the machine")
	itemCmd.PersistentFlags().StringVar(&addParameters, "add-parameters", "", "A JSON string of parameters to add to the machine")

	longDesc := map[string]string{
		"add":      `Add places machines from a pool into the selected pool.  The machines must be unallocated and in Free status.  By default, the default pool is used.`,
		"remove":   `Remove places machines from the selected pool into the default pool.  The machines must be unallocated and in Free status.`,
		"allocate": `Allocate reserves machines in the selected pool.  The machines must be unallocated and in Free status.`,
		"release":  `Release frees machines in the selected pool.  The machines must be allocated and in InUse status.`,
	}

	for _, cmdName := range []string{"add", "remove", "allocate", "release"} {
		cmdAct := &cobra.Command{
			Use:   fmt.Sprintf("%s [id ][filter options a=f(v) style]", cmdName),
			Short: fmt.Sprintf("%s machines to pool", cmdName),
			Long:  longDesc[cmdName],
			Args: func(c *cobra.Command, args []string) error {
				if len(args) == 0 {
					return fmt.Errorf("Must specify a pool id")
				}
				if allMachines && count > 0 {
					return fmt.Errorf("Must choose count or all-machines, but not both")
				}
				if machineList != "" && (count > 0 || allMachines) {
					return fmt.Errorf("machine-list must not be specified with count or all-machines")
				}
				if !allMachines && count == 0 && machineList == "" {
					count = 1
				}
				if minimum > count {
					return fmt.Errorf("count must be greater than minimum")
				}
				if minimum == 0 {
					minimum = count
				}
				if len(args) == 1 {
					return nil
				}
				if strings.Contains(args[1], "=") {
					for i, a := range args {
						if i == 0 {
							continue
						}
						ar := strings.SplitN(a, "=", 2)
						if len(ar) != 2 {
							return fmt.Errorf("Filter argument requires an '=' separator: %s", a)
						}
					}
				}
				return nil
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				parms := map[string]interface{}{
					"pool/all-machines": allMachines,
					"pool/minimum":      minimum,
					"pool/wait-timeout": waitTimeout,
					"pool/count":        count,
					"pool/workflow":     newWorkflow,
				}
				if addProfiles != "" {
					l := []string{}
					for _, pp := range strings.Split(addProfiles, ",") {
						l = append(l, strings.TrimSpace(pp))
					}
					parms["pool/add-profiles"] = l
				}
				if removeProfiles != "" {
					l := []string{}
					for _, pp := range strings.Split(removeProfiles, ",") {
						l = append(l, strings.TrimSpace(pp))
					}
					parms["pool/remove-profiles"] = l
				}
				if removeParameters != "" {
					l := []string{}
					for _, pp := range strings.Split(removeParameters, ",") {
						l = append(l, strings.TrimSpace(pp))
					}
					parms["pool/remove-parameters"] = l
				}
				if addParameters != "" {
					data := map[string]interface{}{}
					if err := bufOrFileDecode(addParameters, &data); err != nil {
						return fmt.Errorf("add-parameters is not a valid JSON or YAML string or file")
					}
					parms["pool/add-parameters"] = data
				}
				if machineList != "" {
					l := []string{}
					for _, pp := range strings.Split(machineList, ",") {
						l = append(l, strings.TrimSpace(pp))
					}
					parms["pool/machine-list"] = l
				}

				id := args[0]
				filters := []string{}
				for i, v := range args {
					if i == 0 {
						continue
					}
					filters = append(filters, v)
				}
				if len(filters) > 0 {
					parms["pool/filter"] = filters
				}

				cmdStr := fmt.Sprintf("%sMachines", cmd.Name())

				pr := []*models.PoolResult{}
				req := Session.Req().Post(parms).UrlFor("pools", id, cmdStr)
				qsParams := []string{}
				if force {
					qsParams = append(qsParams, "force", "true")
				}
				if cmdStr == "addMachines" && sourcePool != "" {
					qsParams = append(qsParams, "source-pool", sourcePool)
				}
				if len(qsParams) > 0 {
					req = req.Params(qsParams...)
				}
				if err := req.Do(&pr); err != nil {
					return err
				}
				return prettyPrint(pr)
			},
		}
		if cmdName == "add" {
			cmdAct.PersistentFlags().StringVar(&sourcePool, "source-pool", "", "The name of the pool to pull machines from")
		}
		itemCmd.AddCommand(cmdAct)
	}
	op.addCommand(itemCmd)
	op.command(app)
}

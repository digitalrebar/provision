package cli

import (
	"fmt"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/spf13/cobra"
)

func (o *ops) params() {
	aggregate := false
	decode := false
	params := ""
	getParams := &cobra.Command{
		Use:   "params [id] [json]",
		Short: fmt.Sprintf("Gets/sets all parameters for the %s", o.singleName),
		Long:  fmt.Sprintf(`A helper function to return all or set all the parameters on the %s`, o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 && len(args) != 2 {
				return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			if len(args) == 1 {
				req := session.Req().UrlFor(o.name, args[0], "params")
				if aggregate {
					req.Params("aggregate", "true")
				}
				if decode {
					req.Params("decode", "true")
				}
				if params != "" {
					req.Params("params", params)
				}
				res := map[string]interface{}{}
				if err := req.Do(&res); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
				return prettyPrint(res)
			}
			val := map[string]interface{}{}
			if err := into(args[1], &val); err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}
			res := map[string]interface{}{}
			if ref == "" {
				if err := session.Req().Post(val).UrlFor(o.name, args[0], "params").Do(&res); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			} else {
				var data map[string]interface{}
				if err := bufOrFileDecode(ref, &data); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				if err := session.Req().ParanoidPatch().PatchObj(data, val).UrlFor(o.name, args[0], "params").Do(&res); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			}
			return prettyPrint(res)
		},
	}
	getParams.Flags().BoolVar(&aggregate, "aggregate", false, "Should return aggregated view")
	getParams.Flags().BoolVar(&decode, "decode", false, "Should return decoded secure params")
	getParams.Flags().StringVar(&params,
		"params",
		"",
		"Should return only the parameters specified as a comma-separated list of parameter names.")
	o.addCommand(getParams)
	getParam := &cobra.Command{
		Use:   "get [id] param [key]",
		Short: fmt.Sprintf("Get a parameter from the %s", o.singleName),
		Long:  fmt.Sprintf(`A helper function to return the value of the parameter on the %s`, o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			key := args[2]
			var res interface{}
			req := session.Req().UrlFor(o.name, uuid, "params", key)
			if aggregate {
				req.Params("aggregate", "true")
			}
			if decode {
				req.Params("decode", "true")
			}
			if err := req.Do(&res); err != nil {
				return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
			}
			return prettyPrint(res)
		},
	}
	getParam.Flags().BoolVar(&aggregate, "aggregate", false, "Should return aggregated view")
	getParam.Flags().BoolVar(&decode, "decode", false, "Should return decoded secure params")
	o.addCommand(getParam)
	o.addCommand(&cobra.Command{
		Use:   "add [id] param [key] to [json blob]",
		Short: fmt.Sprintf("Add the %s param *key* to *blob*", o.name),
		Long:  fmt.Sprintf(`Helper function to add parameters to the %s. Fails is already present.`, o.name),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 5 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			key := args[2]
			newValue := args[4]
			var value interface{}
			err := into(newValue, &value)
			if err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}
			value, err = maybeEncryptParam(key, o.name, uuid, value)
			if err != nil {
				return generateError(err, "Cannot set secure parameter %s", key)
			}

			res := map[string]interface{}{}
			if ref == "" {
				if err := session.Req().UrlFor(o.name, uuid, "params").Do(&res); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			} else {
				if err := bufOrFileDecode(ref, &res); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
			}

			if _, ok := res[key]; ok {
				return fmt.Errorf("Key, %s, already present on %s %s", key, o.singleName, uuid)
			}

			var params interface{}
			path := fmt.Sprintf("/%s", makeJsonPtr(key))
			patch := jsonpatch2.Patch{
				jsonpatch2.Operation{
					Op:    "test",
					Path:  "",
					Value: res,
				},
				jsonpatch2.Operation{
					Op:    "add",
					Path:  path,
					Value: value,
				},
			}
			if err := session.Req().Patch(patch).UrlFor(o.name, uuid, "params").Do(&params); err != nil {
				return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
			}
			return prettyPrint(value)
		},
	})
	o.addCommand(&cobra.Command{
		Use:   "set [id] param [key] to [json blob]",
		Short: fmt.Sprintf("Set the %s param *key* to *blob*", o.name),
		Long:  fmt.Sprintf(`Helper function to update the %s parameters.`, o.name),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 5 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			key := args[2]
			newValue := args[4]
			var value interface{}
			err := into(newValue, &value)
			if err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}
			value, err = maybeEncryptParam(key, o.name, uuid, value)
			if err != nil {
				return generateError(err, "Cannot set secure parameter %s", key)
			}
			var params interface{}
			if ref == "" {
				if err := session.Req().Post(value).UrlFor(o.name, uuid, "params", key).Do(&params); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			} else {
				var data interface{}
				if err := bufOrFileDecode(ref, &data); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				path := fmt.Sprintf("/%s", makeJsonPtr(key))
				patch := jsonpatch2.Patch{
					jsonpatch2.Operation{
						Op:    "test",
						Path:  path,
						Value: data,
					},
					jsonpatch2.Operation{
						Op:    "replace",
						Path:  path,
						Value: value,
					},
				}
				if err := session.Req().Patch(patch).UrlFor(o.name, uuid, "params").Do(&params); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			}
			return prettyPrint(value)
		},
	})
	o.addCommand(&cobra.Command{
		Use:   "remove [id] param [key]",
		Short: fmt.Sprintf("Remove the param *key* from %s", o.name),
		Long:  fmt.Sprintf(`Helper function to update the %s parameters.`, o.name),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			key := args[2]
			var param interface{}
			if ref == "" {
				err := session.Req().Del().UrlFor(o.name, uuid, "params", key).Do(&param)
				if err != nil {
					return generateError(err, "Failed to delete param %v: %v", key, uuid)
				}
			} else {
				var data interface{}
				if err := bufOrFileDecode(ref, &data); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				path := fmt.Sprintf("/%s", makeJsonPtr(key))
				patch := jsonpatch2.Patch{
					jsonpatch2.Operation{
						Op:    "test",
						Path:  path,
						Value: data,
					},
					jsonpatch2.Operation{
						Op:   "remove",
						Path: path,
					},
				}
				if err := session.Req().Patch(patch).UrlFor(o.name, uuid, "params").Do(&param); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			}
			return prettyPrint(param)
		},
	})
}

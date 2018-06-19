package cli

import (
	"fmt"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func (o *ops) meta() {
	getMeta := func(objType, id string) (map[string]string, error) {
		req := session.Req().UrlFor("meta", objType, id)
		res := map[string]string{}
		if err := req.Do(&res); err != nil {
			return nil, generateError(err, "Failed to fetch meta %v: %v", o.singleName, id)
		}
		return res, nil
	}
	mc := &cobra.Command{
		Use:   "meta [id]",
		Short: fmt.Sprintf("Gets metadata for the %s", o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			res, err := getMeta(o.name, id)
			if err != nil {
				return err
			}
			return prettyPrint(res)
		},
	}
	o.addCommand(mc)
	mc.AddCommand(&cobra.Command{
		Use:   "get [id] [key]",
		Short: fmt.Sprintf("Get a specific metadata item from %s", o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			key := args[1]
			res, err := getMeta(o.name, id)
			if err != nil {
				return err
			}
			if val, ok := res[key]; ok {
				return prettyPrint(val)
			}
			return fmt.Errorf("No such metadata item %s", key)
		},
	})
	mc.AddCommand(&cobra.Command{
		Use:   "add [id] key [key] val [val]",
		Short: fmt.Sprintf("Atomically add [key]:[val] to the metadata on [%s]:[id]", o.name),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			key := args[2]
			newValue := args[4]
			var value string
			err := into(newValue, &value)
			if err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}
			res := models.Meta{}
			if ref == "" {
				res, err = getMeta(o.name, id)
				if err != nil {
					return err
				}
			} else {
				if err := bufOrFileDecode(ref, &res); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
			}
			if _, ok := res[key]; ok {
				return fmt.Errorf("Key, %s, already present on %s %s", key, o.singleName, id)
			}

			var meta interface{}
			path := fmt.Sprintf("/%s", makeJsonPtr(key))
			patch := jsonpatch2.Patch{
				{Op: "test", Path: "", Value: res},
				{Op: "add", Path: path, Value: value},
			}
			if err := session.Req().Patch(patch).UrlFor("meta", o.name, id).Do(&meta); err != nil {
				return generateError(err, "Failed to fetch meta %v: %v", o.singleName, id)
			}
			return prettyPrint(value)
		},
	})
	mc.AddCommand(&cobra.Command{
		Use:   "set [id] key [key] to [val]",
		Short: fmt.Sprintf("Set metadata [key]:[val] on [%s]:[id]", o.name),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			key := args[2]
			newValue := args[4]
			var value string
			err := into(newValue, &value)
			if err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}
			patch := jsonpatch2.Patch{}
			path := fmt.Sprintf("/%s", makeJsonPtr(key))
			if ref == "" {
				md, err := getMeta(o.name, id)
				if err != nil {
					return err
				}
				if meta, ok := md[key]; !ok {
					patch = append(patch, jsonpatch2.Operation{Op: "test", Path: "", Value: md})
				} else {
					patch = append(patch, jsonpatch2.Operation{Op: "test", Path: path, Value: meta})
				}
			} else {
				var meta string
				if err := bufOrFileDecode(ref, &meta); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				patch = append(patch, jsonpatch2.Operation{Op: "test", Path: path, Value: meta})
			}
			patch = append(patch, jsonpatch2.Operation{Op: "add", Path: path, Value: value})
			if err := session.Req().Patch(patch).UrlFor("meta", o.name, id).Do(nil); err != nil {
				return generateError(err, "Failed to fetch meta %v: %v", o.singleName, id)
			}
			return prettyPrint(value)
		},
	})
	mc.AddCommand(&cobra.Command{
		Use:   "remove [id] key [key]",
		Short: fmt.Sprintf("Remove the meta [key] from [%s]:[id]", o.name),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			key := args[2]
			path := fmt.Sprintf("/%s", makeJsonPtr(key))
			var patch jsonpatch2.Patch
			if ref == "" {
				patch = jsonpatch2.Patch{
					{
						Op:   "remove",
						Path: path,
					},
				}
			} else {
				var data string
				if err := bufOrFileDecode(ref, &data); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				patch = jsonpatch2.Patch{
					{
						Op:    "test",
						Path:  path,
						Value: data,
					},
					{
						Op:   "remove",
						Path: path,
					},
				}
			}
			if err := session.Req().Patch(patch).UrlFor("meta", o.name, id).Do(nil); err != nil {
				return generateError(err, "Failed to fetch meta %v: %v", o.singleName, id)
			}
			return nil
		},
	})
}

package cli

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"

	utils2 "github.com/VictorLowther/jsonpatch2/utils"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func (o *ops) params() {
	compose := false
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
				req := Session.Req().UrlFor(o.name, args[0], "params")
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
				return fmt.Errorf("Unable to unmarshal input stream: %v", err)
			}
			res := map[string]interface{}{}
			if ref == "" {
				if err := Session.Req().Post(val).UrlFor(o.name, args[0], "params").Do(&res); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			} else {
				var data map[string]interface{}
				if err := bufOrFileDecode(ref, &data); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				if err := Session.Req().ParanoidPatch().PatchObj(data, val).UrlFor(o.name, args[0], "params").Do(&res); err != nil {
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
			req := Session.Req().UrlFor(o.name, uuid, "params", key)
			if aggregate {
				req.Params("aggregate", "true")
			}
			if decode {
				req.Params("decode", "true")
			}
			if compose {
				req.Params("compose", "true")
			}
			if err := req.Do(&res); err != nil {
				return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
			}
			return prettyPrint(res)
		},
	}
	getParam.Flags().BoolVar(&aggregate, "aggregate", false, "Should return aggregated view")
	getParam.Flags().BoolVar(&compose, "compose", false, "Should merge map and array objects together")
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
				return fmt.Errorf("Unable to unmarshal input stream: %v", err)
			}
			value, err = maybeEncryptParam(key, o.name, uuid, value)
			if err != nil {
				return generateError(err, "Cannot set secure parameter %s", key)
			}

			res := map[string]interface{}{}
			if ref == "" {
				if err := Session.Req().UrlFor(o.name, uuid, "params").Do(&res); err != nil {
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
			path := fmt.Sprintf("/%s", makeJSONPtr(key))
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
			if err := Session.Req().Patch(patch).UrlFor(o.name, uuid, "params").Do(&params); err != nil {
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
				return fmt.Errorf("Unable to unmarshal input stream: %v", err)
			}
			value, err = maybeEncryptParam(key, o.name, uuid, value)
			if err != nil {
				return generateError(err, "Cannot set secure parameter %s", key)
			}
			var params interface{}
			if ref == "" {
				if err := Session.Req().Post(value).UrlFor(o.name, uuid, "params", key).Do(&params); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			} else {
				var data interface{}
				if err := bufOrFileDecode(ref, &data); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				path := fmt.Sprintf("/%s", makeJSONPtr(key))
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
				if err := Session.Req().Patch(patch).UrlFor(o.name, uuid, "params").Do(&params); err != nil {
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
				err := Session.Req().Del().UrlFor(o.name, uuid, "params", key).Do(&param)
				if err != nil {
					return generateError(err, "Failed to delete param %v: %v", key, uuid)
				}
			} else {
				var data interface{}
				if err := bufOrFileDecode(ref, &data); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				path := fmt.Sprintf("/%s", makeJSONPtr(key))
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
				if err := Session.Req().Patch(patch).UrlFor(o.name, uuid, "params").Do(&param); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			}
			return prettyPrint(param)
		},
	})
	o.addCommand(&cobra.Command{
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
			name := args[0]
			key := "bootenv-customize"
			var res interface{}
			req := Session.Req().UrlFor(o.name, name, "params", key)
			req.Params("aggregate", "true")
			req.Params("decode", "true")
			req.Params("compose", "true")
			if err := req.Do(&res); err != nil {
				return generateError(err, "Failed to fetch params %v: %v", o.singleName, name)
			}
			if res == nil {
				return fmt.Errorf("%s %s does not require an iso, parameter not set", o.singleName, name)
			}

			bootEnvs := map[string]*models.BootEnv{}
			if jerr := utils2.Remarshal(res, &bootEnvs); jerr != nil {
				return generateError(jerr, "Failed to fetch %v: %v", o.singleName, args[0])
			}

			isoFiles := map[string]string{}
			for _, bootEnv := range bootEnvs {
				if bootEnv.OS.IsoFile != "" {
					isoFiles[bootEnv.OS.IsoFile] = bootEnv.OS.IsoUrl
				}
				for _, archInfo := range bootEnv.OS.SupportedArchitectures {
					if archInfo.IsoFile != "" {
						isoFiles[archInfo.IsoFile] = archInfo.IsoUrl
					}
				}
			}

			if len(isoFiles) == 0 {
				return fmt.Errorf("%s %s does not require an iso", o.singleName, name)
			}
			isos, err := Session.ListBlobs("isos")
			if err != nil {
				return fmt.Errorf("%s %s Unable to determine what ISO files are already present", o.singleName, name)
			}
			for _, iso := range isos {
				if _, ok := isoFiles[iso]; ok {
					delete(isoFiles, iso)
				}
			}
			if len(isoFiles) == 0 {
				log.Printf("%s %s already has all required ISO files", o.singleName, name)
				return nil
			}
			for isoFile, isoUrl := range isoFiles {
				if isoUrl == "" {
					log.Printf("Unable to automatically download iso for %s %s, skipping", o.singleName, name)
					continue
				}
				tr := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}
				if downloadProxy != "" {
					proxyURL, err := url.Parse(downloadProxy)
					if err == nil {
						tr.Proxy = http.ProxyURL(proxyURL)
					}
				}
				client := &http.Client{Transport: tr}
				isoUrl, _ = signRackNUrl(isoUrl)
				isoDlResp, err := client.Get(isoUrl)
				if err != nil {
					log.Printf("Unable to connect to %s: %v: Skipping", isoUrl, err)
					continue
				}
				if isoDlResp.StatusCode >= 300 {
					isoDlResp.Body.Close()
					log.Printf("Unable to initiate download of %s: %s: Skipping", isoUrl, isoDlResp.Status)
					continue
				}
				func() {
					defer isoDlResp.Body.Close()
					if info, err := Session.PostBlob(isoDlResp.Body, "isos", isoFile); err != nil {
						log.Printf("%v", generateError(err, "Error uploading %s", isoUrl))
					} else {
						log.Printf("%v", prettyPrint(info))
					}
				}()
			}
			return nil
		},
	})
}

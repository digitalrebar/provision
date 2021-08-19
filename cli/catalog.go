package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/url"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/Masterminds/semver"

	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func getLocalCatalog() (res *models.Content, err error) {
	req := Session.Req().List("files")
	req.Params("path", "rebar-catalog/rackn-catalog")
	data := []interface{}{}
	err = req.Do(&data)
	if err != nil {
		return
	}

	if len(data) == 0 {
		err = fmt.Errorf("Failed to find local catalog")
		return
	}

	vs := make([]*semver.Version, len(data))
	vmap := map[string]string{}
	for i, obj := range data {
		r := obj.(string)
		v, verr := semver.NewVersion(strings.TrimSuffix(r, ".json"))
		if verr != nil {
			err = verr
			return
		}
		vs[i] = v
		vmap[v.String()] = r
	}
	sort.Sort(sort.Reverse(semver.Collection(vs)))

	var buf bytes.Buffer
	path := fmt.Sprintf("rebar-catalog/rackn-catalog/%s", vmap[vs[0].String()])
	fmt.Printf("Using catalog: %s\n", path)
	if gerr := Session.GetBlob(&buf, "files", path); gerr != nil {
		err = fmt.Errorf("Failed to fetch %v: %v: %v", "files", path, gerr)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &res)
	return
}

func fetchCatalog() (res *models.Content, err error) {
	buf := []byte{}
	buf, err = bufOrFile(catalog)
	if err == nil {
		err = json.Unmarshal(buf, &res)
	}
	if err != nil {
		err = fmt.Errorf("Error fetching catalog: %v", err)
	}
	return
}

func itemsFromCatalog(cat *models.Content, name string) map[string]*models.CatalogItem {
	res := map[string]*models.CatalogItem{}
	for k, v := range cat.Sections["catalog_items"] {
		item := &models.CatalogItem{}
		if err := models.Remarshal(v, &item); err != nil {
			continue
		}
		if name == "" || name == item.Name {
			res[k] = item
		}
	}
	return res
}

func oneItem(cat *models.Content, name, version string) *models.CatalogItem {
	items := itemsFromCatalog(cat, name)
	for _, v := range items {
		if v.Version == version {
			return v
		}
	}
	return nil
}

func installItem(catalog *models.Content, name, version, arch, tgtos string, replaceWritable bool, inflight map[string]struct{}) error {
	inflight[name] = struct{}{}
	if name == "BasicStore" {
		return nil
	}
	item := oneItem(catalog, name, version)
	if item == nil {
		return fmt.Errorf("%s version %s not in catalog", name, version)
	}
	src, err := urlOrFileAsReadCloser(item.DownloadUrl(arch, tgtos))
	if src != nil {
		defer src.Close()
	}
	if err != nil {
		return fmt.Errorf("Unable to contact source URL for %s: %v", item.Name, err)
	}
	switch item.ContentType {
	case "ContentPackage":
		content := &models.Content{}
		if err := json.NewDecoder(src).Decode(&content); err != nil {
			return fmt.Errorf("Error downloading content bundle %s: %v", item.Name, err)
		}

		prereqs := strings.Split(content.Meta.Prerequisites, ",")
		for _, p := range prereqs {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			pversion := version
			if pversion != "tip" && pversion != "stable" {
				pversion = "stable"
			}
			parts := strings.Split(p, ":")
			pname := strings.TrimSpace(parts[0])

			if _, err := Session.GetContentItem(pname); err == nil {
				inflight[pname] = struct{}{}
				continue
			}

			if err := installItem(catalog, pname, pversion, arch, tgtos, replaceWritable, inflight); err != nil {
				return err
			}
		}
		return doReplaceContent(content, "", replaceWritable)
	case "PluginProvider":
		res := &models.PluginProviderUploadInfo{}
		req := Session.Req().Post(src).UrlFor("plugin_providers", item.Name)
		if replaceWritable {
			req = req.Params("replaceWritable", "true")
		}
		// TODO: One day handle prereqs.  Save to local file, mark executable, get contents, check prereqs
		if err := req.Do(res); err != nil {
			return err
		}
		return prettyPrint(res)
	case "DRP":
		if info, err := Session.PostBlob(src, "system", "upgrade"); err != nil {
			return generateError(err, "Failed to post upgrade of DRP")
		} else {
			return prettyPrint(info)
		}
	case "DRPUX":
		if info, err := Session.PostBlobExplode(src, true, "files", "ux", "drp-ux.zip"); err != nil {
			return generateError(err, "Failed to post upgrade of DRP")
		} else {
			return prettyPrint(info)
		}
	default:
		return fmt.Errorf("Don't know how to install %s of type %s yet", item.Name, item.ContentType)
	}
}

func catalogCommands() *cobra.Command {

	type catItem struct {
		Type     string
		Versions []string
	}

	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "Access commands related to catalog manipulation",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show the contents of the current catalog",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			catalog, err := fetchCatalog()
			if err != nil {
				return err
			}
			return prettyPrint(catalog)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "items",
		Short: "Show the items available in the catalog",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			catalog, err := fetchCatalog()
			if err != nil {
				return err
			}

			items := map[string]catItem{}
			for _, v := range itemsFromCatalog(catalog, "") {
				item := &models.CatalogItem{}
				if err := models.Remarshal(v, &item); err != nil {
					continue
				}
				if _, ok := items[item.Name]; !ok {
					items[item.Name] = catItem{Type: item.ContentType, Versions: []string{item.Version}}
				} else {
					cat := items[item.Name]
					cat.Versions = append(cat.Versions, item.Version)
					items[item.Name] = cat
				}
			}
			for k := range items {
				sort.Strings(items[k].Versions)
			}
			return prettyPrint(items)
		},
	})
	itemCmd := &cobra.Command{
		Use:   "item",
		Short: "Commands to act on individual catalog items",
	}
	var arch, tgtos, version string
	itemCmd.PersistentFlags().StringVar(&arch, "arch", runtime.GOARCH, "Architecture of the item to work with when downloading a plugin provider")
	itemCmd.PersistentFlags().StringVar(&tgtos, "os", runtime.GOOS, "OS of the item to work with when downloading a plugin provider")
	itemCmd.PersistentFlags().StringVar(&version, "version", "stable", "Version of the item to work with")
	itemCmd.AddCommand(&cobra.Command{
		Use:   "download [item] (to [file])",
		Short: "Downloads [item] to [file]",
		Long: `Downloads the specified item to the specified file
If to [file] is not specified, it will be downloaded into current directory
and wind up in a file with the same name as the item + the default file extension for the file type.
`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 && len(args) != 3 {
				return fmt.Errorf("item download requires 1 or 2 arguments")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			catalog, err := fetchCatalog()
			if err != nil {
				return err
			}
			item := oneItem(catalog, args[0], version)
			if item == nil {
				return fmt.Errorf("%s version %s not in catalog", args[0], version)
			}
			target := item.FileName()
			if len(args) == 3 {
				target = args[2]
			}
			mode := os.FileMode(0644)
			if item.ContentType == "PluginProvider" {
				mode = 0755
			}
			src, err := urlOrFileAsReadCloser(item.DownloadUrl(arch, tgtos))
			if src != nil {
				defer src.Close()
			}
			if err != nil {
				return fmt.Errorf("Unable to contact source URL for %s: %v", item.Name, err)
			}
			fi, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
			if err != nil {
				return fmt.Errorf("Unable to create %s: %v", target, err)
			}
			defer fi.Close()
			_, err = io.Copy(fi, src)
			return err
		},
	})
	replaceWritable := false
	install := &cobra.Command{
		Use:               "install [item]",
		Short:             "Installs [item] from the catalog on the current dr-provision endpoint",
		PersistentPreRunE: ppr,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("item install requires 1 argument")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			catalog, err := fetchCatalog()
			if err != nil {
				return err
			}
			info, err := Session.Info()
			if err != nil {
				return fmt.Errorf("Unable to fetch session information to determine endpoint arch and OS")
			}
			arch = info.Arch
			tgtos = info.Os
			err = installItem(catalog, args[0], version, arch, tgtos, replaceWritable, map[string]struct{}{})
			if err != nil {
				return err
			}
			return nil
		},
	}
	install.Flags().BoolVar(&replaceWritable, "replace-writable", false, "Replace identically named writable objects")
	// Flag deprecated due to standardizing on all hyphenated form for persistent flags.
	install.Flags().BoolVar(&replaceWritable, "replaceWritable", false, "Replace identically named writable objects")
	install.Flags().MarkHidden("replaceWritable")
	install.Flags().MarkDeprecated("replaceWritable", "please use --replace-writable")
	itemCmd.AddCommand(install)
	itemCmd.AddCommand(&cobra.Command{
		Use:   "show [item]",
		Short: "Shows available versions for [item]",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("item show requires 1 argument")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			catalog, err := fetchCatalog()
			if err != nil {
				return err
			}

			items := map[string]catItem{}
			for _, v := range itemsFromCatalog(catalog, args[0]) {
				item := &models.CatalogItem{}
				if err := models.Remarshal(v, &item); err != nil {
					continue
				}
				if _, ok := items[item.Name]; !ok {
					items[item.Name] = catItem{Type: item.ContentType, Versions: []string{item.Version}}
				} else {
					cat := items[item.Name]
					cat.Versions = append(cat.Versions, item.Version)
					items[item.Name] = cat
				}
			}
			if len(items) == 0 {
				return fmt.Errorf("No item named %s in the catalog", args[0])
			}
			for k := range items {
				sort.Strings(items[k].Versions)
			}
			return prettyPrint(items[args[0]])
		},
	})
	cmd.AddCommand(itemCmd)

	minVersion := ""
	tip := false
	concurrency := 1
	updateCmd := &cobra.Command{
		Use:   "updateLocal",
		Short: "Update the local catalog from the upstream catalog",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			if concurrency > 20 {
				return fmt.Errorf("Invalid value for concurrency: %d. Max allowed is 20", concurrency)
			}

			srcCatalog, err := fetchCatalog()
			if err != nil {
				return err
			}

			localCatalog, err := getLocalCatalog()
			if err != nil {
				return err
			}

			srcItems := itemsFromCatalog(srcCatalog, "")
			localItems := itemsFromCatalog(localCatalog, "")

			requireStable := false
			if minVersion == "stable" {
				requireStable = true
			}

			var mv *semver.Version
			if !requireStable {
				var verr error
				mv, verr = semver.NewVersion(minVersion)
				if verr != nil {
					return fmt.Errorf("Invalid version: %s, %v", minVersion, verr)
				}
			}

			// We want to be able to process `concurrency` number of updates at a time
			srcItemsKeys := make([]string, 0)
			srcItemsKeyChunks := make([][]string, 0)
			for k, _ := range srcItems {
				srcItemsKeys = append(srcItemsKeys, k)
			}
			// Collect chunks of keys to process
			for i := 0; i < len(srcItemsKeys); i += concurrency {
				end := i + concurrency
				// Avoid out of range
				if end > len(srcItemsKeys) {
					end = len(srcItemsKeys)
				}
				srcItemsKeyChunks = append(srcItemsKeyChunks, srcItemsKeys[i:end])
			}

			var g errgroup.Group
			for _, srcItemsKeys := range srcItemsKeyChunks {
				// Assign temporary variable that is local to this iteration
				srcItemsKeysIteration := srcItemsKeys
				g.Go(func() error {
					for _, k := range srcItemsKeysIteration {
						if minVersion == "" {
							// We want to get everything here except tip (unless tip is set to true in which case we include tip)
							if !tip && srcItems[k].Tip {
								continue
							}
						} else if requireStable {
							// Only get things that are stable or tip if tip is set to true
							if !(srcItems[k].Version == "stable" || (tip && srcItems[k].Tip)) {
								continue
							}
						} else {
							// Only get things that aren't stable
							if srcItems[k].Version == "stable" {
								continue
							}
							// Only get versions greater than the input version excluding tip (unless tip is set to true in which case we include tip)
							if mv != nil {
								if o, oerr := semver.NewVersion(srcItems[k].ActualVersion); oerr != nil || mv.Compare(o) > 0 || (!tip && srcItems[k].Tip) {
									continue
								}
							}
						}

						// Get things that aren't in local
						if nv, ok := localItems[k]; !ok || nv.ActualVersion != srcItems[k].ActualVersion {
							parts := map[string]string{}
							i := strings.Index(srcItems[k].Source, "/rebar-catalog/")
							switch srcItems[k].ContentType {
							case "DRP":
								for arch := range srcItems[k].Shasum256 {
									if arch == "any/any" {
										parts[srcItems[k].Source], _ = url.QueryUnescape(srcItems[k].Source[i+1:])
									} else {
										archValue := strings.Split(arch, "/")[0]
										osValue := strings.Split(arch, "/")[1]
										ts := strings.ReplaceAll(srcItems[k].Source, ".zip", "."+archValue+"."+osValue+".zip")
										qs, _ := url.QueryUnescape(srcItems[k].Source[i+1:])
										td := strings.ReplaceAll(qs, ".zip", "."+archValue+"."+osValue+".zip")
										parts[ts] = td
									}
								}
							case "PluginProvider", "DRPCLI":
								for arch := range srcItems[k].Shasum256 {
									ts := fmt.Sprintf("%s/%s/%s", srcItems[k].Source, arch, srcItems[k].Name)
									qs, _ := url.QueryUnescape(srcItems[k].Source[i+1:])
									td := fmt.Sprintf("%s/%s/%s", qs, arch, srcItems[k].Name)
									parts[ts] = td
								}
							default:
								parts[srcItems[k].Source], _ = url.QueryUnescape(srcItems[k].Source[i+1:])
							}

							for s, d := range parts {
								fmt.Printf("Downloading %s to store at %s\n", s, d)
								data, err := urlOrFileAsReadCloser(s)
								if err != nil {
									return fmt.Errorf("Error opening src file %s: %v", s, err)
								}
								func() {
									defer data.Close()
									_, err = Session.PostBlobExplode(data, false, "files", d)
								}()
								if err != nil {
									return generateError(err, "Failed to post %v: %v", "files", d)
								}
							}
						}
					}
					return nil
				})
				if err := g.Wait(); err != nil {
					return err
				}
			}
			return nil
		},
	}
	updateCmd.PersistentFlags().StringVar(&minVersion, "version", "", "Minimum version of the items. If set to 'stable' will only get stable entries.")
	updateCmd.PersistentFlags().BoolVar(&tip, "tip", false, "Include tip versions of the packages")
	updateCmd.PersistentFlags().IntVar(&concurrency, "concurrency", 1, "Number of concurrent download options")
	cmd.AddCommand(updateCmd)

	// Start of create stuff
	var pkgVer string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a custom catalog for Digital Rebar",
		RunE: func(cmd *cobra.Command, args []string) error {
			cat, err := fetchCatalog()
			if err != nil {
				return err
			}
			newMap := map[string]interface{}{}
			for k, v := range cat.Sections["catalog_items"] {
				item := &models.CatalogItem{}
				if err := models.Remarshal(v, &item); err != nil {
					continue
				}
				if item.Version == pkgVer {
					newMap[k] = item
				}
			}
			cat.Sections["catalog_items"] = newMap
			return prettyPrint(cat)
		},
	}
	createCmd.Flags().StringVarP(&pkgVer, "pkg-version", "", "", "pkg-version tip|stable (required)")
	createCmd.MarkFlagRequired("pkg-version")
	cmd.AddCommand(createCmd)

	return cmd
}

func init() {
	addRegistrar(func(c *cobra.Command) { c.AddCommand(catalogCommands()) })
}

func installOne() {

}

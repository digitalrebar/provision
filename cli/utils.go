package cli

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/olekukonko/tablewriter"
)

var encodeJSONPtr = strings.NewReplacer("~", "~0", "/", "~1")

// String translates a pointerSegment into a regular string, encoding it as we go.
func makeJSONPtr(s string) string {
	return encodeJSONPtr.Replace(string(s))
}

func bufOrFileDecode(ref string, data interface{}) (err error) {
	buf, terr := bufOrStdin(ref)
	if terr != nil {
		err = fmt.Errorf("Unable to process reference object: %v", terr)
		return
	}
	err = api.DecodeYaml(buf, &data)
	if err != nil {
		err = fmt.Errorf("Unable to unmarshal reference object: %v", err)
		return
	}
	return
}

func getCatalogSource(nv string) (string, error) {
	// XXX: Query self first?  One day.

	catalogData, err := bufOrFile(catalog)
	if err != nil {
		return "", err
	}

	clayer := &models.Content{}
	if err := api.DecodeYaml(catalogData, clayer); err != nil {
		return "", err
	}

	var elem interface{}
	for k, cobj := range clayer.Sections["catalog_items"] {
		if k == nv {
			elem = cobj
			break
		}
	}
	if elem == nil {
		return "", fmt.Errorf("Catalog item: %s not found", nv)
	}

	ci := &models.CatalogItem{}
	if err := utils.Remarshal(elem, &ci); err != nil {
		return "", fmt.Errorf("Catalog item: %s can not be remarshaled: %v", nv, err)
	}
	if ci.Source == "" {
		return "", fmt.Errorf("Catalog item: %s does not have a source: %v", nv, ci)
	}
	return ci.DownloadUrl(runtime.GOOS, runtime.GOARCH), nil
}

func urlOrFileAsReadCloser(src string) (io.ReadCloser, error) {
	if s, err := os.Lstat(src); err == nil && s.Mode().IsRegular() {
		fi, err := os.Open(src)
		if err != nil {
			return nil, fmt.Errorf("Error opening %s: %v", src, err)
		}
		return fi, nil
	}
	if strings.HasPrefix(src, "catalog:") {
		var err error
		src, err = getCatalogSource(strings.TrimPrefix(src, "catalog:"))
		if err != nil {
			return nil, err
		}
	}
	if u, err := url.Parse(src); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
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
		res, err := client.Get(src)
		if err != nil {
			return nil, err
		}
		return res.Body, nil
	} else if err == nil && u.Scheme == "file" {
		return nil, fmt.Errorf("file:// scheme not supported")
	}
	return nil, fmt.Errorf("Must specify a file or url")
}

func bufOrFile(src string) ([]byte, error) {
	if s, err := os.Lstat(src); err == nil && s.Mode().IsRegular() {
		return ioutil.ReadFile(src)
	}
	if strings.HasPrefix(src, "catalog:") {
		var err error
		src, err = getCatalogSource(strings.TrimPrefix(src, "catalog:"))
		if err != nil {
			return nil, err
		}
	}
	if u, err := url.Parse(src); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
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
		res, err := client.Get(src)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		return []byte(body), err
	} else if err == nil && u.Scheme == "file" {
		return nil, fmt.Errorf("file:// scheme not supported")
	}
	return []byte(src), nil
}

func bufOrStdin(src string) ([]byte, error) {
	if src == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return bufOrFile(src)
}

func into(src string, res interface{}) error {
	buf, err := bufOrStdin(src)
	if err != nil {
		return fmt.Errorf("Error reading from stdin: %v", err)
	}
	return api.DecodeYaml(buf, &res)
}

func safeMergeJSON(src interface{}, toMerge []byte) ([]byte, error) {
	sv := reflect.ValueOf(src)
	for sv.Kind() == reflect.Ptr || sv.Kind() == reflect.Interface {
		sv = sv.Elem()
	}
	toMergeObj := make(map[string]interface{})
	if err := json.Unmarshal(toMerge, &toMergeObj); err != nil {
		return nil, err
	}
	targetObj := map[string]interface{}{}
	if err := utils.Remarshal(src, &targetObj); err != nil {
		return nil, err
	}
	outObj, ok := utils.Merge(targetObj, toMergeObj).(map[string]interface{})
	if !ok {
		return nil, errors.New("Cannot happen in safeMergeJSON")
	}
	if sv.Kind() == reflect.Struct {
		finalObj := map[string]interface{}{}
		for i := 0; i < sv.NumField(); i++ {
			vf := sv.Field(i)
			if !vf.CanSet() {
				continue
			}
			tf := sv.Type().Field(i)
			mapField := tf.Name
			if tag, ok := tf.Tag.Lookup(`json`); ok {
				tagVals := strings.Split(tag, `,`)
				if tagVals[0] == "-" {
					continue
				}
				if tagVals[0] != "" {
					mapField = tagVals[0]
				}
			}
			if v, ok := outObj[mapField]; ok {
				finalObj[mapField] = v
			}
		}
		return json.Marshal(finalObj)
	}
	// For Raw!!
	return json.Marshal(outObj)
}

func mergeInto(src models.Model, changes []byte) (models.Model, error) {
	buf, err := safeMergeJSON(src, changes)
	if err != nil {
		return nil, err
	}
	dest := models.Clone(src)
	return dest, json.Unmarshal(buf, &dest)
}

func mergeFromArgs(src models.Model, changes string) (models.Model, error) {
	// We have to load this and then convert to json to merge safely.
	data := map[string]interface{}{}
	if err := bufOrFileDecode(changes, &data); err != nil {
		return nil, err
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return mergeInto(src, buf)
}

func d(msg string, args ...interface{}) {
	if debug {
		log.Printf(msg, args...)
	}
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func lamePrinter(obj interface{}) []byte {
	isTable := format == "table"

	if slice, ok := obj.([]interface{}); ok {
		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)

		var theFields []string
		data := [][]string{}

		for i, v := range slice {
			if m, ok := v.(map[string]interface{}); ok {
				if i == 0 {
					theFields = strings.Split(printFields, ",")
					if printFields == "" {
						theFields = []string{}
						for k := range m {
							theFields = append(theFields, k)
						}
					}
				}
				row := []string{}
				for _, k := range theFields {
					row = append(row, truncateString(fmt.Sprintf("%v", m[k]), truncateLength))
				}
				data = append(data, row)
			} else {
				if i == 0 {
					theFields = []string{"Index", "Value"}
				}
				data = append(data, []string{fmt.Sprintf("%d", i), truncateString(fmt.Sprintf("%v", obj), truncateLength)})
			}
		}

		if !noHeader {
			table.SetHeader(theFields)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetHeaderLine(isTable)
		}
		table.SetAutoWrapText(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		if !isTable {
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetTablePadding("\t") // pad with tabs
			table.SetBorder(false)
			table.SetNoWhiteSpace(true)
		}
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
		return []byte(tableString.String())
	}
	if m, ok := obj.(map[string]interface{}); ok {
		theFields := strings.Split(printFields, ",")
		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)

		if !noHeader {
			table.SetHeader([]string{"Field", "Value"})
			table.SetHeaderLine(isTable)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		}
		table.SetAutoWrapText(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		if !isTable {
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetTablePadding("\t") // pad with tabs
			table.SetBorder(false)
			table.SetNoWhiteSpace(true)
		}

		data := [][]string{}

		if printFields != "" {
			for _, k := range theFields {
				data = append(data, []string{k, truncateString(fmt.Sprintf("%v", m[k]), truncateLength)})
			}
		} else {
			index := []string{}
			for k := range m {
				index = append(index, k)
			}
			sort.Strings(index)
			for _, k := range index {
				v := m[k]
				data = append(data, []string{k, truncateString(fmt.Sprintf("%v", v), truncateLength)})
			}
		}

		table.AppendBulk(data) // Add Bulk Data
		table.Render()
		return []byte(tableString.String())
	}

	// Default for everything else
	return []byte(truncateString(fmt.Sprintf("%v", obj), truncateLength))
}

func prettyPrintBuf(o interface{}) (buf []byte, err error) {
	var v interface{}
	if err := utils.Remarshal(o, &v); err != nil {
		return nil, err
	}

	if format == "text" || format == "table" {
		return lamePrinter(v), nil
	}
	return api.Pretty(format, v)
}

func prettyPrint(o interface{}) (err error) {
	if noPretty {
		fmt.Printf("%v", o)
		return nil
	}
	buf, err := prettyPrintBuf(o)
	if err != nil {
		return err
	}
	fmt.Println(string(buf))
	return nil
}

package cli

import (
	"bytes"
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
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-isatty"

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

var rackNHost = regexp.MustCompile(`(?P<bucket>[-\w]+)\.s3[-.](?P<region>[-\w]+)\.amazonaws.com`)
var rackNUrl = regexp.MustCompile(`s3[-.](?P<region>[-\w]+)\.amazonaws.com`)

type CloudiaUrlReq struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
	Object string `json:"object"`
}

type CloudiaUrlResp struct {
	Object CloudiaUrlReq `json:"object"`
	Url    string        `json:"url"`
}

func signRackNUrl(oldUrl string) (newUrl string, err error) {
	newUrl = oldUrl
	myURL, perr := url.Parse(oldUrl)
	if perr != nil {
		err = perr
		return
	}

	// Test Time
	values := myURL.Query()
	expire := values.Get("X-Amz-Expires")
	date := values.Get("X-Amz-Date")
	// Check if we already have a signed URL.
	if expire != "" && date != "" {
		d, derr := time.Parse("20060102T150405Z0700", date)
		sec, serr := time.ParseDuration(expire + "s")
		if derr == nil && serr == nil {
			if time.Now().Before(d.Add(sec)) {
				return
			}
		}
	}

	// Should we sign the URL - test the hostname format
	match := rackNHost.FindStringSubmatch(myURL.Host)
	result := map[string]string{}
	for i, name := range rackNHost.SubexpNames() {
		if i != 0 && name != "" && len(match) > i {
			result[name] = match[i]
		}
	}
	bucket, bok := result["bucket"]
	region, rok := result["region"]
	myPath := strings.TrimPrefix(myURL.Path, "/")
	if !bok || !rok {
		// Check the region/bucket form
		match = rackNUrl.FindStringSubmatch(myURL.Host)
		result = map[string]string{}
		for i, name := range rackNUrl.SubexpNames() {
			if i != 0 && name != "" && len(match) > i {
				result[name] = match[i]
			}
		}
		region, rok = result["region"]
		if !rok {
			err = fmt.Errorf("Url is not a RackN URL: %s", oldUrl)
			return
		}
		parts := strings.SplitN(myPath, "/", 2)
		if len(parts) != 2 {
			err = fmt.Errorf("Url is not a RackN URL: %s", oldUrl)
			return
		}
		bucket = parts[0]
		myPath = parts[1]
	}

	if Session == nil {
		err = fmt.Errorf("No session to get signing info")
		return
	}

	// Something to sign.
	license := ""
	if lerr := Session.Req().UrlFor("profiles", "rackn-license", "params", "rackn/license").Do(&license); err != nil {
		err = fmt.Errorf("Failed to get license: %v", lerr)
		return
	}

	reqData := &CloudiaUrlReq{
		Bucket: bucket,
		Region: region,
		Object: myPath,
	}
	tr := &http.Transport{
		MaxIdleConns:    1,
		IdleConnTimeout: 10 * time.Second,
	}
	data, _ := json.Marshal(reqData)
	client := &http.Client{Transport: tr}
	req, rerr := http.NewRequest("POST", "https://cloudia.rackn.io/api/v1/org/presign", bytes.NewBuffer(data))
	if rerr != nil {
		err = fmt.Errorf("Failed to build request for cloudia: %v", rerr)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", license)

	resp, derr := client.Do(req)
	if derr != nil {
		err = fmt.Errorf("Failed to query cloudia: %v", derr)
		return
	}
	defer resp.Body.Close()

	body, berr := ioutil.ReadAll(resp.Body)
	if berr != nil {
		err = fmt.Errorf("Failed to read body: %v", berr)
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Cloudia returned error: %d: %s", resp.StatusCode, string(body))
		return
	}

	respData := &CloudiaUrlResp{}
	if jerr := json.Unmarshal(body, &respData); jerr != nil {
		err = fmt.Errorf("Failed to unmarshal response: %s %v", string(body), jerr)
		return
	}

	newUrl = respData.Url
	return
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
		src, _ = signRackNUrl(src)
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
		src, _ = signRackNUrl(src)
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

		colColors := []tablewriter.Colors{}
		headerColors := []tablewriter.Colors{}
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
					if !noColor {
						for range theFields {
							headerColors = append(headerColors, tablewriter.Color(colorPatterns[4]...))
							colColors = append(colColors, tablewriter.Color(colorPatterns[6]...))
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
					if !noColor {
						headerColors = []tablewriter.Colors{tablewriter.Color(colorPatterns[4]...), tablewriter.Color(colorPatterns[5]...)}
						colColors = []tablewriter.Colors{tablewriter.Color(colorPatterns[6]...), tablewriter.Color(colorPatterns[7]...)}
					}
				}
				data = append(data, []string{fmt.Sprintf("%d", i), truncateString(fmt.Sprintf("%v", obj), truncateLength)})
			}
		}

		if !noHeader {
			table.SetHeader(theFields)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetHeaderLine(isTable)
			if !noColor {
				table.SetHeaderColor(headerColors...)
				table.SetColumnColor(colColors...)
			}
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
			if !noColor {
				table.SetHeaderColor(tablewriter.Color(colorPatterns[4]...), tablewriter.Color(colorPatterns[5]...))
				table.SetColumnColor(tablewriter.Color(colorPatterns[6]...), tablewriter.Color(colorPatterns[7]...))
			}
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

var colorPatterns [][]int

func processColorPatterns() {
	if colorPatterns != nil {
		return
	}

	colorPatterns = [][]int{
		// JSON
		[]int{32},    // String
		[]int{33},    // Bool
		[]int{36},    // Number
		[]int{90},    // Null
		[]int{34, 1}, // Key
		// Table colors
		[]int{35}, // Header
		[]int{92}, // Value
		[]int{32}, // Header2
		[]int{35}, // Value2
	}

	parts := strings.Split(colorString, ";")
	for _, p := range parts {
		subparts := strings.Split(p, "=")
		idx, e := strconv.Atoi(subparts[0])
		if e != nil {
			continue
		}
		if idx < 0 || idx >= len(colorPatterns) {
			continue
		}
		attrs := strings.Split(subparts[1], ",")
		if len(attrs) == 0 {
			continue
		}
		ii := make([]int, len(attrs))
		for i, attr := range attrs {
			ii[i], e = strconv.Atoi(attr)
			if e != nil {
				ii = nil
				break
			}
		}
		if ii != nil {
			colorPatterns[idx] = ii
		}
	}
}

func prettyPrintBuf(o interface{}) (buf []byte, err error) {
	var v interface{}
	if err := utils.Remarshal(o, &v); err != nil {
		return nil, err
	}

	noColor = noColor || os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
	processColorPatterns()

	if format == "text" || format == "table" {
		return lamePrinter(v), nil
	}
	return api.PrettyColor(format, v, !noColor, colorPatterns)
}

func prettyPrint(o interface{}) (err error) {
	var buf []byte
	buf, err = prettyPrintBuf(o)
	if err != nil {
		return
	}
	fmt.Println(string(buf))
	if errHaver, ok := o.(models.Validator); ok && objectErrorsAreFatal {
		err = errHaver.HasError()
	}
	return
}

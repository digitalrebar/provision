package models

import (
	"bytes"
	"encoding/json"
	"log"
	"path"
)

/*
 * CatalogItem tracks a versioned thing in the RackN system
 */

// CatalogItem structure that handles RawModel instead of dealing with
// RawModel which is how DRP is storing it.
type CatalogItem struct {
	Validation
	Access
	// Meta Items
	// Icon        string
	// Color       string
	// Author      string
	// DisplayName string
	// License     string
	// Copyright   string
	// CodeSource  string
	Meta

	Owned

	Id   string
	Type string

	Name          string
	ActualVersion string
	Version       string
	ContentType   string
	Source        string
	Shasum256     map[string]string

	Tip    bool
	HotFix bool
}

func (ci *CatalogItem) Key() string {
	return ci.Id
}

func (ci *CatalogItem) KeyName() string {
	return "Id"
}

func (ci *CatalogItem) Prefix() string {
	return "catalog_items"
}

// Clone the endpoint
func (ci *CatalogItem) Clone() *CatalogItem {
	ci2 := &CatalogItem{}
	buf := bytes.Buffer{}
	enc, dec := json.NewEncoder(&buf), json.NewDecoder(&buf)
	if err := enc.Encode(ci); err != nil {
		log.Panicf("Failed to encode endpoint:%s: %v", ci.Id, err)
	}
	if err := dec.Decode(ci2); err != nil {
		log.Panicf("Failed to decode endpoint:%s: %v", ci.Id, err)
	}
	return ci2
}

func (ci *CatalogItem) Fill() {
	ci.Type = "catalog_items"
	if ci.Meta == nil {
		ci.Meta = Meta{}
	}
	if ci.Errors == nil {
		ci.Errors = []string{}
	}
	if ci.Shasum256 == nil {
		ci.Shasum256 = map[string]string{}
	}
}

// DownloadUrl returns a URL that you can use to download the artifact
// for this catalog item.  If the CatalogItem has a ContentType of
// `PluginProvider`, arch and os must be set appropriately for
// the target binary type, otherwise they can be left blank.
func (ci *CatalogItem) DownloadUrl(arch, os string) string {
	res := ci.Source
	switch ci.ContentType {
	case "PluginProvider":
		res = res + path.Join("/", arch, os, ci.Name)
	case "DRPCLI":
		res = res + path.Join("/", arch, os, ci.Name)
		if os == "windows" {
			res = res + ".exe"
		}
	}
	return res
}

// FileName returns the recommended filename to use when writing this catalog item to disk.
func (ci *CatalogItem) FileName() string {
	switch ci.ContentType {
	case "DRP", "DRPUX":
		return ci.Name + ".zip"
	case "ContentPackage":
		return ci.Name + ".json"
	default:
		return ci.Name
	}
}

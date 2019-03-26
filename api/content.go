package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/digitalrebar/provision/models"
	"github.com/digitalrebar/store"
)

func (c *Client) GetContentSummary() ([]*models.ContentSummary, error) {
	res := []*models.ContentSummary{}
	return res, c.Req().UrlFor("contents").Do(&res)
}

func (c *Client) GetContentItem(name string) (*models.Content, error) {
	res := &models.Content{}
	return res, c.Req().UrlFor("contents", name).Do(&res)
}

func (c *Client) CreateContent(content *models.Content) (*models.ContentSummary, error) {
	res := &models.ContentSummary{}
	return res, c.Req().Post(content).UrlFor("contents").Do(res)
}

func (c *Client) ReplaceContent(content *models.Content) (*models.ContentSummary, error) {
	res := &models.ContentSummary{}
	return res, c.Req().Put(content).UrlFor("contents", content.Meta.Name).Do(res)
}

func (c *Client) DeleteContent(name string) error {
	return c.Req().Del().UrlFor("contents", name).Do(nil)
}

func FindOrFake(src, field string, args map[string]string) string {
	if src != "" {
		filepath := fmt.Sprintf("._%s.meta", field)
		buf, err := ioutil.ReadFile(path.Join(src, filepath))
		if err == nil {
			return string(buf)
		}
	}
	if p, ok := args[field]; !ok {
		s := "Unspecified"
		switch field {
		case "Type":
			// Default Type should be dynamic
			s = "dynamic"
		case "RequiredFeatures", "Version", "Prerequisites":
			// Default RequiredFeatures should be empty string
			s = ""
		}
		return s
	} else {
		return p
	}
}

func (c *Client) BundleContent(src string, dst store.Store, params map[string]string) error {
	if dm, ok := dst.(store.MetaSaver); ok {
		meta := map[string]string{
			"Name":             FindOrFake(src, "Name", params),
			"Version":          FindOrFake(src, "Version", params),
			"Description":      FindOrFake(src, "Description", params),
			"Source":           FindOrFake(src, "Source", params),
			"Documentation":    FindOrFake(src, "Documentation", params),
			"RequiredFeatures": FindOrFake(src, "RequiredFeatures", params),
			"Type":             FindOrFake(src, "Type", params),
			"Color":            FindOrFake(src, "Color", params),
			"Icon":             FindOrFake(src, "Icon", params),
			"Author":           FindOrFake(src, "Author", params),
			"DisplayName":      FindOrFake(src, "DisplayName", params),
			"License":          FindOrFake(src, "License", params),
			"Copyright":        FindOrFake(src, "Copyright", params),
			"CodeSource":       FindOrFake(src, "CodeSource", params),
			"Order":            FindOrFake(src, "Order", params),
			"Tags":             FindOrFake(src, "Tags", params),
			"DocUrl":           FindOrFake(src, "DocUrl", params),
			"Prerequisites":    FindOrFake(src, "Prerequisites", params),
		}
		dm.SetMetaData(meta)
	}

	// for each valid content type, load it
	files, _ := ioutil.ReadDir(src)
	for _, f := range files {
		if !f.IsDir() || (f.IsDir() && strings.HasPrefix(f.Name(), ".")) {
			continue
		}
		prefix := f.Name()

		if _, err := models.New(prefix); err != nil {
			// Skip things we can instantiate
			continue
		}
		sub, err := dst.MakeSub(prefix)
		if err != nil {
			return fmt.Errorf("Cannot make substore %s: %v", prefix, err)
		}
		items, err := ioutil.ReadDir(path.Join(src, prefix))
		if err != nil {
			return fmt.Errorf("Cannot read substore %s: %v", prefix, err)
		}
		for _, fileInfo := range items {
			if fileInfo.IsDir() {
				continue
			}
			itemName := fileInfo.Name()
			item, _ := models.New(prefix)
			buf, err := ioutil.ReadFile(path.Join(src, prefix, itemName))
			if err != nil {
				return fmt.Errorf("Cannot read item %s: %v", path.Join(prefix, itemName), err)
			}
			switch path.Ext(itemName) {
			case ".yaml", ".yml":
				if err := store.YamlCodec.Decode(buf, item); err != nil {
					return fmt.Errorf("Cannot parse item %s: %v", path.Join(prefix, itemName), err)
				}
			case ".json":
				if err := store.JsonCodec.Decode(buf, item); err != nil {
					return fmt.Errorf("Cannot parse item %s: %v", path.Join(prefix, itemName), err)
				}
			default:
				if tmpl, ok := item.(*models.Template); ok && prefix == "templates" {
					tmpl.ID = itemName
					tmpl.Contents = string(buf)
				} else {
					return fmt.Errorf("No idea how to decode %s into %s", itemName, item.Prefix())
				}
			}
			if err := sub.Save(item.Key(), item); err != nil {
				return fmt.Errorf("Failed to save %s:%s: %v", item.Prefix(), item.Key(), err)
			}
		}
	}
	return nil
}

func writeMetaFile(dst, field, data string) error {
	if data == "" {
		return nil
	}
	fname := fmt.Sprintf("._%s.meta", field)
	return ioutil.WriteFile(path.Join(dst, fname), []byte(data), 0640)
}

func (c *Client) UnbundleContent(content store.Store, dst string) error {
	if err := os.MkdirAll(dst, 0750); err != nil {
		return err
	}
	if cm, ok := content.(store.MetaSaver); ok {
		meta := cm.MetaData()
		for k, v := range meta {
			if err := writeMetaFile(dst, k, v); err != nil {
				return err
			}
		}
	}
	for prefix, sub := range content.Subs() {
		if err := os.MkdirAll(path.Join(dst, prefix), 0750); err != nil {
			return err
		}
		_, err := models.New(prefix)
		if err != nil {
			return fmt.Errorf("Store contains model of type %s the we don't know about", prefix)
		}
		keys, err := sub.Keys()
		if err != nil {
			return fmt.Errorf("Failed to retrieve keys for substore %s: %v", prefix, err)
		}
		codec := content.GetCodec()
		for _, key := range keys {
			item, _ := models.New(prefix)
			if err := sub.Load(key, item); err != nil {
				return fmt.Errorf("Failed to load %s:%s: %v", prefix, key, err)
			}
			var buf []byte
			var fname string
			switch prefix {
			case "templates":
				fname = strings.Replace(key, "/", ".", -1)
				buf = []byte(item.(*models.Template).Contents)
			default:
				fname = strings.Replace(key, "/", ".", -1) + codec.Ext()
				buf, err = codec.Encode(item)
				if err != nil {
					return fmt.Errorf("Failed to encode %s:%s: %v", prefix, key, err)
				}
			}
			if err := ioutil.WriteFile(path.Join(dst, prefix, fname), buf, 0640); err != nil {
				return fmt.Errorf("Failed to save %s:%s: %v", prefix, key, err)
			}
		}
	}
	return nil
}

package models

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateInfo holds information on the templates in the boot
// environment that will be expanded into files.
//
// swagger:model
type TemplateInfo struct {
	// Name of the template
	//
	// required: true
	Name string
	// A text/template that specifies how to create
	// the final path the template should be
	// written to.
	//
	// required: true
	Path string
	// The ID of the template that should be expanded.  Either
	// this or Contents should be set
	//
	// required: false
	ID string
	// The contents that should be used when this template needs
	// to be expanded.  Either this or ID should be set.
	//
	// required: false
	Contents string
	// Metadata for the TemplateInfo.  This can be used by the job running
	// system and the bootenvs to handle OS, arch, and firmware differences.
	//
	// required: false
	Meta     map[string]string
	pathTmpl *template.Template
}

func (ti *TemplateInfo) Id() string {
	if ti.ID == "" {
		return ti.Name
	}
	return ti.ID
}

func (ti *TemplateInfo) SanityCheck(idx int, e ErrorAdder, missingPathOK bool) {
	if ti.Name == "" {
		e.Errorf("Template[%d] is missing a Name", idx)
	}
	if !missingPathOK {
		if ti.Path == "" {
			e.Errorf("Template[%d] is missing a Path", idx)
		} else if _, err := template.New(ti.Name).Parse(ti.Path); err != nil {
			e.Errorf("Template[%d] Path is not a valid text/template: %v", idx, err)
		}
	}
	if ti.Contents == "" && ti.ID == "" {
		e.Errorf("Template[%d] must have either an ID or Contents set", idx)
	}
	if ti.Contents != "" && ti.ID != "" {
		e.Errorf("Template[%d] has both an ID and Contents", idx)
	}
	if ti.Meta == nil {
		ti.Meta = map[string]string{}
	}
}

func (ti *TemplateInfo) PathTemplate() *template.Template {
	return ti.pathTmpl
}

func MergeTemplates(root *template.Template, tmpls []TemplateInfo, e ErrorAdder) *template.Template {
	var res *template.Template
	var err error
	if root == nil {
		res = template.New("")
	} else {
		res, err = root.Clone()
	}
	if err != nil {
		e.Errorf("Error cloning root: %v", err)
		return nil
	}
	buf := &bytes.Buffer{}
	for i := range tmpls {
		ti := &tmpls[i]
		if ti.Name == "" {
			e.Errorf("Templates[%d] has no Name", i)
			continue
		}
		if ti.Path != "" {
			pathTmpl, err := template.New(ti.Name).Parse(ti.Path)
			if err != nil {
				e.Errorf("Error compiling path template %s (%s): %v",
					ti.Name,
					ti.Path,
					err)
				continue
			} else {
				ti.pathTmpl = pathTmpl.Option("missingkey=error")
			}
		}
		if ti.ID != "" {
			if res.Lookup(ti.ID) == nil {
				e.Errorf("Templates[%d]: No common template for %s", i, ti.ID)
			}
			continue
		}
		if ti.Contents == "" {
			e.Errorf("Templates[%d] has both an empty ID and contents", i)
		}
		fmt.Fprintf(buf, `{{define "%s"}}%s{{end}}\n`, ti.Name, ti.Contents)
	}
	_, err = res.Parse(buf.String())
	if err != nil {
		e.Errorf("Error parsing inline templates: %v", err)
	}
	return res
}

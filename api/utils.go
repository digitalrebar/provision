package api

import (
	"fmt"

	prettyjson "github.com/hokaccha/go-prettyjson"

	"github.com/fatih/color"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/ghodss/yaml"
)

type Decoder interface {
	Decode(interface{}) error
}

type Encoder interface {
	Encode(interface{}) error
}

// DecodeYaml is a helper function for dealing with user input -- when
// accepting input from the user, we want to treat both YAML and JSON
// as first-class citizens.  The YAML library we use makes that easier
// by using the json struct tags for all marshalling and unmarshalling
// purposes.
//
// Note that the REST API does not use YAML as a wire protocol, so
// this function should never be used to decode data coming from the
// provision service.
func DecodeYaml(buf []byte, ref interface{}) error {
	return models.DecodeYaml(buf, ref)
}

// Pretty marshals object according to the the fmt, in whatever
// passed for "pretty" according to fmt.
func Pretty(f string, obj interface{}) ([]byte, error) {
	return PrettyColor(f, obj, false, nil)
}

// PrettyColor marshals object according to the the fmt, in whatever
// passed for "pretty" according to fmt.  If useColor = true, then
// try to colorize output
func PrettyColor(f string, obj interface{}, useColor bool, colors [][]int) ([]byte, error) {
	switch f {
	case "json":
		f := prettyjson.NewFormatter()
		f.StringColor = color.New(color.FgGreen)
		f.BoolColor = color.New(color.FgYellow)
		f.NumberColor = color.New(color.FgCyan)
		f.NullColor = color.New(color.FgHiBlack)
		f.KeyColor = color.New(color.FgBlue, color.Bold)
		if colors != nil {
			if len(colors) > 0 && colors[0] != nil {
				attrs := make([]color.Attribute, len(colors[0]))
				for i, v := range colors[0] {
					attrs[i] = color.Attribute(v)
				}
				f.StringColor = color.New(attrs...)
			}
			if len(colors) > 1 && colors[1] != nil {
				attrs := make([]color.Attribute, len(colors[1]))
				for i, v := range colors[1] {
					attrs[i] = color.Attribute(v)
				}
				f.BoolColor = color.New(attrs...)
			}
			if len(colors) > 2 && colors[2] != nil {
				attrs := make([]color.Attribute, len(colors[2]))
				for i, v := range colors[2] {
					attrs[i] = color.Attribute(v)
				}
				f.NumberColor = color.New(attrs...)
			}
			if len(colors) > 3 && colors[3] != nil {
				attrs := make([]color.Attribute, len(colors[3]))
				for i, v := range colors[3] {
					attrs[i] = color.Attribute(v)
				}
				f.NullColor = color.New(attrs...)
			}
			if len(colors) > 4 && colors[4] != nil {
				attrs := make([]color.Attribute, len(colors[4]))
				for i, v := range colors[4] {
					attrs[i] = color.Attribute(v)
				}
				f.KeyColor = color.New(attrs...)
			}
		}
		f.DisabledColor = !useColor
		return f.Marshal(obj)
	case "yaml":
		return yaml.Marshal(obj)
	case "go":
		buf, err := Pretty("yaml", obj)
		if err == nil {
			return []byte(fmt.Sprintf("package main\n\nvar contentYamlString = `\n%s\n`\n", string(buf))), nil
		}
		return nil, err
	default:
		return nil, fmt.Errorf("Unknown pretty format %s", f)
	}
}

// GenPatch generates a JSON patch that will transform source into target.
// The generated patch will have all the applicable test clauses.
func GenPatch(source, target interface{}, paranoid bool) (jsonpatch2.Patch, error) {
	return models.GenPatch(source, target, paranoid)
}

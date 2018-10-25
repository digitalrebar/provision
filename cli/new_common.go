package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

type ops struct {
	name          string
	singleName    string
	example       func() models.Model
	mustPut       bool
	noCreate      bool
	noUpdate      bool
	noDestroy     bool
	noWait        bool
	extraCommands []*cobra.Command
	actionName    string
}

func maybeEncryptParam(param string,
	prefix, key string,
	val interface{}) (interface{}, error) {
	p := &models.Param{}
	if err := session.FillModel(p, param); err != nil {
		return val, nil
	}
	if !p.Secure {
		return val, nil
	}
	k := []byte{}
	if err := session.Req().UrlFor(prefix, key, "pubkey").Do(&k); err != nil {
		return nil, err
	}

	sv := &models.SecureData{}
	return sv, sv.Marshal(k, val)
}

func (o *ops) refOrFill(key string) (data models.Model, err error) {
	data = o.example()

	if ref == "" {
		if err = session.FillModel(data, key); err != nil {
			return
		}
	} else {
		err = bufOrFileDecode(ref, &data)
	}
	return
}

func (o *ops) addCommand(c *cobra.Command) {
	o.extraCommands = append(o.extraCommands, c)
}

func (o *ops) command(app *cobra.Command) {
	res := &cobra.Command{
		Use:   o.name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", o.name),
	}
	if o.name == "extended" {
		res.PersistentFlags().StringVarP(&o.name,
			"ldata", "l", "",
			"object type for extended data commands")
	}
	if o.example != nil {
		ref := o.example()
		if _, ok := ref.(models.BootEnver); ok {
			o.bootenv()
		}
		if _, ok := ref.(models.Paramer); ok {
			o.params()
		}
		if _, ok := ref.(models.Profiler); ok {
			o.profiles()
		}
		if _, ok := ref.(models.Tasker); ok {
			o.tasks()
		}
		if _, ok := ref.(models.Actor); ok {
			o.actions()
		}
		if _, ok := ref.(models.MetaHaver); ok {
			o.meta()
		}
		res.AddCommand(o.commands()...)
	}
	app.AddCommand(res)
}

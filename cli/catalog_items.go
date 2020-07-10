package cli

import (
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerCatalogItem)
}

func registerCatalogItem(app *cobra.Command) {
	op := &ops{
		name:       "catalog_items",
		singleName: "catalog_item",
		example:    func() models.Model { return &models.CatalogItem{} },
	}
	op.command(app)
}

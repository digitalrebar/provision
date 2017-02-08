package restapi

import (
	"crypto/tls"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	swag "github.com/go-openapi/swag"
	graceful "github.com/tylerb/graceful"

	provisioner "github.com/galthaus/swagger-test/provisioner"

	"github.com/galthaus/swagger-test/restapi/operations"
	"github.com/galthaus/swagger-test/restapi/operations/bootenvs"
	"github.com/galthaus/swagger-test/restapi/operations/files"
	"github.com/galthaus/swagger-test/restapi/operations/isos"
	"github.com/galthaus/swagger-test/restapi/operations/machines"
	"github.com/galthaus/swagger-test/restapi/operations/templates"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../swagger.yaml

func configureFlags(api *operations.RocketSkatesAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "provisioner",
			LongDescription:  "Provisioner Options",
			Options:          &provisioner.ProvOpts,
		},
	}
}

func configureAPI(api *operations.RocketSkatesAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// GREG: One thing I haven't figured out yet. - how to get that.
	provisioner.InitializeProvisioner(8091)
	api.Logger = provisioner.Logger.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.BinConsumer = runtime.ByteStreamConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.BinProducer = runtime.ByteStreamProducer()

	api.BootenvsDeleteBootenvHandler = bootenvs.DeleteBootenvHandlerFunc(func(params bootenvs.DeleteBootenvParams) middleware.Responder {
		return middleware.NotImplemented("operation bootenvs.DeleteBootenv has not yet been implemented")
	})
	api.BootenvsGetBootenvHandler = bootenvs.GetBootenvHandlerFunc(func(params bootenvs.GetBootenvParams) middleware.Responder {
		return middleware.NotImplemented("operation bootenvs.GetBootenv has not yet been implemented")
	})
	api.BootenvsListBootenvsHandler = bootenvs.ListBootenvsHandlerFunc(func(params bootenvs.ListBootenvsParams) middleware.Responder {
		return middleware.NotImplemented("operation bootenvs.ListBootenvs has not yet been implemented")
	})
	api.BootenvsPatchBootenvHandler = bootenvs.PatchBootenvHandlerFunc(func(params bootenvs.PatchBootenvParams) middleware.Responder {
		return middleware.NotImplemented("operation bootenvs.PatchBootenv has not yet been implemented")
	})
	api.BootenvsPostBootenvHandler = bootenvs.PostBootenvHandlerFunc(func(params bootenvs.PostBootenvParams) middleware.Responder {
		return middleware.NotImplemented("operation bootenvs.PostBootenv has not yet been implemented")
	})
	api.BootenvsPutBootenvHandler = bootenvs.PutBootenvHandlerFunc(func(params bootenvs.PutBootenvParams) middleware.Responder {
		return middleware.NotImplemented("operation bootenvs.PutBootenv has not yet been implemented")
	})

	api.FilesListFilesHandler = files.ListFilesHandlerFunc(provisioner.ListFiles)
	api.FilesGetFileHandler = files.GetFileHandlerFunc(provisioner.GetFile)
	api.FilesPostFileHandler = files.PostFileHandlerFunc(provisioner.UploadFile)
	api.FilesDeleteFileHandler = files.DeleteFileHandlerFunc(provisioner.DeleteFile)

	api.IsosListIsosHandler = isos.ListIsosHandlerFunc(provisioner.ListIsos)
	api.IsosGetIsoHandler = isos.GetIsoHandlerFunc(provisioner.GetIso)
	api.IsosPostIsoHandler = isos.PostIsoHandlerFunc(provisioner.UploadIso)
	api.IsosDeleteIsoHandler = isos.DeleteIsoHandlerFunc(provisioner.DeleteIso)

	api.TemplatesListTemplatesHandler = templates.ListTemplatesHandlerFunc(provisioner.TemplateList)
	api.TemplatesGetTemplateHandler = templates.GetTemplateHandlerFunc(provisioner.TemplateGet)
	api.TemplatesPostTemplateHandler = templates.PostTemplateHandlerFunc(provisioner.TemplatePost)
	api.TemplatesReplaceTemplateHandler = templates.ReplaceTemplateHandlerFunc(provisioner.TemplateReplace)
	api.TemplatesPutTemplateHandler = templates.PutTemplateHandlerFunc(provisioner.TemplatePut)
	api.TemplatesPatchTemplateHandler = templates.PatchTemplateHandlerFunc(provisioner.TemplatePatch)
	api.TemplatesDeleteTemplateHandler = templates.DeleteTemplateHandlerFunc(provisioner.TemplateDelete)

	api.MachinesListMachinesHandler = machines.ListMachinesHandlerFunc(provisioner.MachineList)
	api.MachinesPostMachineHandler = machines.PostMachineHandlerFunc(provisioner.MachinePost)
	api.MachinesGetMachineHandler = machines.GetMachineHandlerFunc(provisioner.MachineGet)
	api.MachinesPutMachineHandler = machines.PutMachineHandlerFunc(provisioner.MachinePut)
	api.MachinesPatchMachineHandler = machines.PatchMachineHandlerFunc(provisioner.MachinePatch)
	api.MachinesDeleteMachineHandler = machines.DeleteMachineHandlerFunc(provisioner.MachineDelete)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.

	// GREG: Do cert.Server config here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	// Serve the swagger UI
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/swagger-ui" || r.URL.Path == "/api/help" {
			http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui/dist"))).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

package plugin

import (
	"github.com/digitalrebar/logger"
	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
)

// PluginStop is an optional interface that your plugin can implement
// to provide custom behaviour whenever the plugin is stopped.
// Stop will be called when a plugin needs to stop operating.
// If implement the PluginStop interface, it will be called
// before the default stop action takes place.
//
// Stop takes one argument, a logger.
type PluginStop interface {
	Stop(logger.Logger)
}

// PluginConfig is a mandatory interface that your plugin must implement.
// Config will be called with the fully expanded Params on the plugin
// object whenever this instance of the Plugin is started or whenever
// those Params change.
//
// Config takes three arguments: a logger, an API client with superuser rights,
// and a map of all the params on the Plugin being configured.  It should
// return a non-nil models.Error if the Config call fails for any reason.
type PluginConfig interface {
	Config(logger.Logger, *api.Client, map[string]interface{}) *models.Error
}

// PluginEventSelector is an optional interface that your plugin can
// implement to specify what events the plugin is interested in receiving
// from dr-provision.  If this interface is implemented, then the HasPublish
// field in the PluginProvider definition must be false.
//
// SelectEvents returns a slice of strings that define the events the
// Plugin wishes to receive.
type PluginEventSelecter interface {
	SelectEvents() []string
}

// PluginPublisher is an optional interface that your plugin can implement
// if it is interested in receiving events from dr-provision.  There are a
// couple of things to be aware of when implementing a Publish method:
//
// 1. If you are implementing a Publish method, you should also implement a
//    SelectEvents method, and set the HasPublish flag on your PluginProvider
//    definition to false. This will allow dr-provision to only send you the
//    specific events you are interested in, and it will prevent your plugin
//    from being able to bottleneck (or even deadlock) dr-provision.
//
// 2. If you choose to not implement a SelectEvents method, the HasPublish
//    flag on your PluginProvider definition must be set to true, and your
//    Publish method will receive all the events dr-provision emits
//    synchronously.  It is therefore your responsibility handle taking action
//    on the events in such a way that you do not cause a deadlock or a
//    performance bottleneck.
//
// Publish takes a logger and the event that was recieved, and returns
// a non-nil models.Error if there was an error handling the event.
type PluginPublisher interface {
	Publish(logger.Logger, *models.Event) *models.Error
}

// PluginActor is an optional interface that your plugin should implement
// if you plan on handling Actions.  If the PluginProvider definition has
// a non-empty list of AvailableActions, then the Action method must
// be available and able to handle all of the Actions in that list.
//
// Action takes a logger and a fully-filled out Action, and returns
// the results of that action along with a non-nil models.Error
// if an error occurred while performing the action.
type PluginActor interface {
	Action(logger.Logger, *models.Action) (interface{}, *models.Error)
}

// PluginValidator is an optional interface that your plugin can implement
// if it needs to check that it can run in the environment it was executed in.
// Validate is a good method to implement to test for other executables, etc.
// that the Plugin may rely on to operate.
//
// Validate takes a logger and an API client with superuser permissions,
// and returns the results of validating the environment and a non-nil
// models.Error if the plugin cannot be used in the current environment.
type PluginValidator interface {
	Validate(logger.Logger, *api.Client) (interface{}, *models.Error)
}

// PluginUnpacker is an optional interface that your plugin can implement
// if it needs to unpack additional assets into the static file space.
//
// Unpack takes a logger and the location on the local filesystem any
// embedded assets should be unpacked into.  It returns an error if
// there was a problem unpacking assets.
type PluginUnpacker interface {
	Unpack(logger.Logger, string) error
}

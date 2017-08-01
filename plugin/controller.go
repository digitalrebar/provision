package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/backend/index"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
)

// Plugin Provider describes the available functions that could be
// instantiated by a plugin.
// swagger:model
type PluginProvider struct {
	Name    string
	Version string

	HasPublish       bool
	AvailableActions []*AvailableAction

	RequiredParams []string
	OptionalParams []string

	// Ensure that these are in the system.
	Parameters []*backend.Param

	path string
}

// swagger:model
type PluginProviderUploadInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

type RunningPlugin struct {
	Plugin   *backend.Plugin
	Provider *PluginProvider
	Client   *PluginClient
}

type PluginController struct {
	logger             *log.Logger
	lock               sync.Mutex
	AvailableProviders map[string]*PluginProvider
	runningPlugins     map[string]*RunningPlugin
	dataTracker        *backend.DataTracker
	pluginDir          string
	watcher            *fsnotify.Watcher
	done               chan bool
	finished           chan bool
	events             chan *backend.Event
	publishers         *backend.Publishers
	MachineActions     *MachineActions
	apiPort            int
}

func InitPluginController(pluginDir string, dt *backend.DataTracker, logger *log.Logger, pubs *backend.Publishers, apiPort int) (pc *PluginController, err error) {
	pc = &PluginController{pluginDir: pluginDir, dataTracker: dt, publishers: pubs, logger: logger,
		AvailableProviders: make(map[string]*PluginProvider, 0), apiPort: apiPort,
		runningPlugins: make(map[string]*RunningPlugin, 0)}

	pc.MachineActions = NewMachineActions()
	pubs.Add(pc)

	pc.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}

	err = pc.watcher.Add(pc.pluginDir)
	if err != nil {
		return
	}

	pc.done = make(chan bool)
	pc.finished = make(chan bool)
	pc.events = make(chan *backend.Event, 1000)

	go func() {
		done := false
		for !done {
			select {
			case event := <-pc.watcher.Events:
				// Skip events on the parent directory.
				if event.Name == pc.pluginDir {
					continue
				}
				// Skip downloading files
				if strings.HasSuffix(event.Name, ".part") {
					continue
				}
				arr := strings.Split(event.Name, "/")
				file := arr[len(arr)-1]
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					pc.lock.Lock()
					pc.removePluginProvider(file)
					pc.lock.Unlock()
				} else if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Chmod == fsnotify.Chmod {
					pc.lock.Lock()
					pc.importPluginProvider(file)
					pc.lock.Unlock()
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					pc.logger.Printf("Rename file: %s %v\n", event.Name, event)
				} else {
					pc.logger.Println("Unhandled file event:", event.Name)
				}
			case event := <-pc.events:
				if event.Action == "create" {
					pc.lock.Lock()
					ref := dt.NewPlugin()
					d, unlocker := dt.LockEnts(ref.Locks("get")...)
					ref2 := d(ref.Prefix()).Find(event.Key)
					// May be deleted before we get here.
					if ref2 != nil {
						pc.startPlugin(d, ref2.(*backend.Plugin))
					}
					unlocker()
					pc.lock.Unlock()
				} else if event.Action == "save" {
					pc.lock.Lock()
					ref := dt.NewPlugin()
					d, unlocker := dt.LockEnts(ref.Locks("get")...)
					ref2 := d(ref.Prefix()).Find(event.Key)
					// May be deleted before we get here.
					if ref2 != nil {
						pc.restartPlugin(d, ref2.(*backend.Plugin))
					}
					unlocker()
					pc.lock.Unlock()
				} else if event.Action == "update" {
					pc.lock.Lock()
					ref := dt.NewPlugin()
					d, unlocker := dt.LockEnts(ref.Locks("get")...)
					// May be deleted before we get here.
					ref2 := d(ref.Prefix()).Find(event.Key)
					if ref2 != nil {
						pc.restartPlugin(d, ref2.(*backend.Plugin))
					}
					unlocker()
					pc.lock.Unlock()
				} else if event.Action == "delete" {
					pc.lock.Lock()
					pc.stopPlugin(event.Object.(*backend.Plugin))
					pc.lock.Unlock()
				} else {
					pc.logger.Println("internal event:", event)
				}
			case err := <-pc.watcher.Errors:
				pc.logger.Println("error:", err)
			case <-pc.done:
				done = true
			}
		}
		pc.finished <- true
	}()

	pc.lock.Lock()
	defer pc.lock.Unlock()

	// Walk plugin dir contents with lock
	files, err := ioutil.ReadDir(pc.pluginDir)
	if err != nil {
		return
	}
	for _, f := range files {
		err = pc.importPluginProvider(f.Name())
		if err != nil {
			return
		}

	}

	return
}

func (pc *PluginController) walkPlugins(provider string) (err error) {
	// Walk all plugin objects from dt.
	var idx *index.Index
	ref := pc.dataTracker.NewPlugin()
	d, unlocker := pc.dataTracker.LockEnts(ref.Locks("get")...)
	defer unlocker()
	idx, err = index.All([]index.Filter{index.Native()}...)(&d(ref.Prefix()).Index)
	if err != nil {
		return
	}
	arr := idx.Items()
	for _, res := range arr {
		plugin := res.(*backend.Plugin)
		if plugin.Provider == provider {
			pc.startPlugin(d, plugin)
		}
	}
	return
}

func (pc *PluginController) Shutdown(ctx context.Context) error {
	pc.done <- true
	<-pc.finished
	return pc.watcher.Close()
}

func (pc *PluginController) Publish(e *backend.Event) error {
	if e.Type != "plugins" {
		return nil
	}
	pc.events <- e
	return nil
}

// This never gets unloaded.
func (pc *PluginController) Reserve() error {
	return nil
}
func (pc *PluginController) Release() {}
func (pc *PluginController) Unload()  {}

func (pc *PluginController) GetPluginProvider(name string) *PluginProvider {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	if pp, ok := pc.AvailableProviders[name]; !ok {
		return nil
	} else {
		return pp
	}
}

func (pc *PluginController) GetPluginProviders() []*PluginProvider {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	// get the list of keys and sort them
	keys := []string{}
	for key := range pc.AvailableProviders {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	answer := []*PluginProvider{}
	for _, key := range keys {
		answer = append(answer, pc.AvailableProviders[key])
	}
	return answer
}

func (pc *PluginController) startPlugin(d backend.Stores, plugin *backend.Plugin) {
	pc.logger.Printf("Starting plugin: %s(%s)\n", plugin.Name, plugin.Provider)
	if _, ok := pc.runningPlugins[plugin.Name]; ok {
		pc.logger.Printf("Already started plugin: %s(%s)\n", plugin.Name, plugin.Provider)
	}
	pp, ok := pc.AvailableProviders[plugin.Provider]
	if ok {
		errors := []string{}

		for _, parmName := range pp.RequiredParams {
			obj, ok := plugin.Params[parmName]
			if !ok {
				errors = append(errors, fmt.Sprintf("Missing required parameter: %s", parmName))
			} else {
				pobj := d("params").Find(parmName)
				if pobj != nil {
					rp := pobj.(*backend.Param)

					if ev := rp.Validate(obj); ev != nil {
						errors = append(errors, ev.Error())
					}
				}
			}
		}
		for _, parmName := range pp.OptionalParams {
			obj, ok := plugin.Params[parmName]
			if ok {
				pobj := d("params").Find(parmName)
				if pobj != nil {
					rp := pobj.(*backend.Param)

					if ev := rp.Validate(obj); ev != nil {
						errors = append(errors, ev.Error())
					}
				}
			}
		}

		if len(errors) == 0 {
			thingee, err := NewPluginClient(plugin.Name, pc.logger, pc.apiPort, pp.path, plugin.Params)
			if err == nil {
				rp := &RunningPlugin{Plugin: plugin, Client: thingee, Provider: pp}
				if pp.HasPublish {
					pc.publishers.Add(thingee)
				}
				for _, aa := range pp.AvailableActions {
					aa.Provider = pp.Name
					aa.plugin = rp
					pc.MachineActions.Add(aa)
				}
				pc.runningPlugins[plugin.Name] = rp
			} else {
				errors = append(errors, err.Error())
			}
		}

		if len(plugin.Errors) != len(errors) {
			plugin.Errors = errors
			pc.dataTracker.Update(d, plugin, nil)
		}
		pc.publishers.Publish("plugin", "started", plugin.Name, plugin)
		pc.logger.Printf("Starting plugin: %s(%s) complete\n", plugin.Name, plugin.Provider)
	} else {
		pc.logger.Printf("Starting plugin: %s(%s) missing provider\n", plugin.Name, plugin.Provider)
		if plugin.Errors == nil || len(plugin.Errors) == 0 {
			plugin.Errors = []string{fmt.Sprintf("Missing Plugin Provider: %s", plugin.Provider)}
			pc.dataTracker.Update(d, plugin, nil)
		}
	}
}

func (pc *PluginController) stopPlugin(plugin *backend.Plugin) {
	rp, ok := pc.runningPlugins[plugin.Name]
	if ok {
		pc.logger.Printf("Stopping plugin: %s(%s)\n", plugin.Name, plugin.Provider)
		delete(pc.runningPlugins, plugin.Name)

		if rp.Provider.HasPublish {
			pc.publishers.Remove(rp.Client)
		}
		for _, aa := range rp.Provider.AvailableActions {
			pc.MachineActions.Remove(aa)
		}
		rp.Client.Stop()
		pc.logger.Printf("Stoping plugin: %s(%s) complete\n", plugin.Name, plugin.Provider)
		pc.publishers.Publish("plugin", "stopped", plugin.Name, plugin)
	}
}

func (pc *PluginController) restartPlugin(d backend.Stores, plugin *backend.Plugin) {
	pc.logger.Printf("Restarting plugin: %s(%s)\n", plugin.Name, plugin.Provider)
	pc.stopPlugin(plugin)
	pc.startPlugin(d, plugin)
	pc.logger.Printf("Restarting plugin: %s(%s) complete\n", plugin.Name, plugin.Provider)
}

// Try to add to available - Must lock before calling
func (pc *PluginController) importPluginProvider(provider string) error {
	pc.logger.Printf("Importing plugin provider: %s\n", provider)
	out, err := exec.Command(pc.pluginDir+"/"+provider, "define").Output()
	if err != nil {
		pc.logger.Printf("Skipping %s because %s\n", provider, err)
	} else {
		var pp PluginProvider
		err = json.Unmarshal(out, &pp)
		if err != nil {
			pc.logger.Printf("Skipping %s because of bad json: %s\n%s\n", provider, err, out)
		} else {

			skip := false
			for _, p := range pp.Parameters {
				err := p.BeforeSave()
				if err != nil {
					pc.logger.Printf("Skipping %s because of bad required scheme: %s %s\n", pp.Name, p.Name, err)
					skip = true
				} else {
					// Attempt create if it doesn't exist already.
					ref := pc.dataTracker.NewParam()
					d, unlocker := pc.dataTracker.LockEnts(ref.Locks("create")...)
					ref2 := d(ref.Prefix()).Find(p.Name)
					if ref2 == nil {
						if _, err := pc.dataTracker.Create(d, p, nil); err != nil {
							pc.logger.Printf("Skipping %s because parameter could not be created: %s %s\n", pp.Name, p.Name, err)
							skip = true
						}
					}
					unlocker()
				}
			}

			if !skip {
				if _, ok := pc.AvailableProviders[pp.Name]; !ok {
					pc.logger.Printf("Adding plugin provider: %s\n", pp.Name)
					pc.AvailableProviders[pp.Name] = &pp
					pp.path = pc.pluginDir + "/" + provider
					for _, aa := range pp.AvailableActions {
						aa.Provider = pp.Name
					}
					pc.publishers.Publish("plugin_provider", "create", pp.Name, pp)
					return pc.walkPlugins(provider)
				} else {
					pc.logger.Printf("Already exists plugin provider: %s\n", pp.Name)
				}
			}
		}
	}
	return nil
}

// Try to stop using plugins and remove available - Must lock before calling
func (pc *PluginController) removePluginProvider(provider string) {
	var name string
	for _, pp := range pc.AvailableProviders {
		if provider == pp.Name {
			name = pp.Name
			break
		}
	}
	if name != "" {
		remove := []*backend.Plugin{}
		for _, rp := range pc.runningPlugins {
			if rp.Provider.Name == name {
				remove = append(remove, rp.Plugin)
			}
		}
		for _, p := range remove {
			ref := pc.dataTracker.NewPlugin()
			d, unlocker := pc.dataTracker.LockEnts(ref.Locks("get")...)
			pc.stopPlugin(p)
			ref2 := d(ref.Prefix()).Find(p.Name)
			myPP := ref2.(*backend.Plugin)
			myPP.Errors = []string{fmt.Sprintf("Missing Plugin Provider: %s", provider)}
			pc.dataTracker.Update(d, myPP, nil)
			unlocker()
		}

		pc.logger.Printf("Removing plugin provider: %s\n", name)
		pc.publishers.Publish("plugin_provider", "delete", name, pc.AvailableProviders[name])
		delete(pc.AvailableProviders, name)
	}
}

func (pc *PluginController) UploadPlugin(c *gin.Context, name string) (*PluginProviderUploadInfo, *backend.Error) {
	if c.Request.Header.Get(`Content-Type`) != `application/octet-stream` {
		return nil, backend.NewError("API ERROR", http.StatusUnsupportedMediaType,
			fmt.Sprintf("upload: plugin_provider %s must have content-type application/octet-stream", name))
	}
	if c.Request.Body == nil {
		return nil, backend.NewError("API ERROR", http.StatusBadRequest,
			fmt.Sprintf("upload: Unable to upload %s: missing body", name))
	}

	ppTmpName := path.Join(pc.pluginDir, fmt.Sprintf(`.%s.part`, path.Base(name)))
	ppName := path.Join(pc.pluginDir, path.Base(name))
	if _, err := os.Open(ppTmpName); err == nil {
		return nil, backend.NewError("API ERROR", http.StatusConflict, fmt.Sprintf("upload: plugin_provider %s already uploading", name))
	}
	tgt, err := os.Create(ppTmpName)
	if err != nil {
		return nil, backend.NewError("API ERROR", http.StatusConflict, fmt.Sprintf("upload: Unable to upload %s: %v", name, err))
	}

	copied, err := io.Copy(tgt, c.Request.Body)
	if err != nil {
		os.Remove(ppTmpName)
		return nil, backend.NewError("API ERROR",
			http.StatusInsufficientStorage, fmt.Sprintf("upload: Failed to upload %s: %v", name, err))
	}
	if c.Request.ContentLength > 0 && copied != c.Request.ContentLength {
		os.Remove(ppTmpName)
		return nil, backend.NewError("API ERROR", http.StatusBadRequest,
			fmt.Sprintf("upload: Failed to upload entire file %s: %d bytes expected, %d bytes received", name, c.Request.ContentLength, copied))
	}
	os.Remove(ppName)
	os.Rename(ppTmpName, ppName)
	os.Chmod(ppName, 0700)
	return &PluginProviderUploadInfo{Path: name, Size: copied}, nil
}

func (pc *PluginController) RemovePlugin(name string) error {
	pluginProviderName := path.Join(pc.pluginDir, path.Base(name))
	return os.Remove(pluginProviderName)
}

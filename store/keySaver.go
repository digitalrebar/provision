package store

import "fmt"

// KeySaver is the interface that should be satisfied by anything that
// wants to use the generic CRUD functions that the Store package
// provides in the keySaver.go file. The methods that should be
// implemented are:
//
// Prefix() should return something that can be used to disambiguate
// object types.  It is only used for error generating methods in the
// store package.
//
// Key() should return a per-Prefix() unique value.
//
// New() should return a fresh uninitialized copy of whatever implements KeySaver.
//
// Backend() should return the store being used to back the object.
type KeySaver interface {
	Prefix() string
	Key() string
	KeyName() string
	New() KeySaver
}

//LoadHooker is the interface that things can satisfy if they want to
//run a hook against an object each time it is loaded from a
//Backend().  OnLoad() will be called after the object has been loaded
//from the Backend().
type LoadHooker interface {
	KeySaver
	OnLoad() error
}

// ChangeHooker is the interface that things can satisfy if they want
// to run a hook against an object that is being Update()'ed.
// OnChange() will be called against the to-be-saved object before it
// is saved to the Backend() with a copy of the object as it currently
// exists in the Backend().  If OnChange() returns a non-nil error, the
// object will no0t be saved.
type ChangeHooker interface {
	KeySaver
	OnChange(KeySaver) error
}

// BeforeDeleteHooker is the interface that things can satisfy if they
// want to test things before an object is removed from the backing
// Backend(). BeforeDelete() will be called before the object is removed
// from the store, and if it returns a non-nil error the object wil
// not be deleted.
type BeforeDeleteHooker interface {
	KeySaver
	BeforeDelete() error
}

// AfterDeleteHooker is the interface things can satisfy if they want
// to perform an action after removal from their backing Backend().
// AfterDelete() is called after the object has been removed from the
// Backend().
type AfterDeleteHooker interface {
	KeySaver
	AfterDelete()
}

// CreateHooker is the interface things can satisfy if they want to
// perfrom an action before a new object is saved to its backing
// store.  OnCreate() will be called just after Create() verifies that
// no object with the same Key() is in the Backend(), and if it returns a
// non-nil error the object will not be saved.
type CreateHooker interface {
	KeySaver
	OnCreate() error
}

// BeforeSaveHooker is the interface things can satisfy if they want
// to perform an action just before the object is saved to the
// Backend().  It is called after any OnCreate() or OnChange() hooks.
// If BeforeSave returns an error, the object will not be saved to the
// Backend().
type BeforeSaveHooker interface {
	KeySaver
	BeforeSave() error
}

// SaveCleanHooker is the interface things can satisfy if they have
// parts that should not be persisted.  It is called immediately
// before the object is saved, and after BeforeSaveHooker finshes.
// It should return a copy of itself that has been "cleaned up".
type SaveCleanHooker interface {
	KeySaver
	SaveClean() KeySaver
}

// AfterSaveHooker is the interface things can satisfy if they want to
// perfrom an action after an object is saved to the Backend().
// AfterSave() will be called after the object has been sucessfully
// saved.
type AfterSaveHooker interface {
	KeySaver
	AfterSave()
}

// When loading an object, the store should
// set the current readonly state on the object.
type ReadOnlySetter interface {
	SetReadOnly(bool)
}

// When loading an object, the store should
// set the current owner name.
type BundleSetter interface {
	SetBundle(string)
}

func load(s Store, k KeySaver, key string, runhook bool) (bool, error) {
	err := s.Load(key, k)
	if err != nil {
		return false, err
	}
	if h, ok := k.(LoadHooker); runhook && ok {
		return true, h.OnLoad()
	}
	return true, nil
}

// List returns a slice of KeySavers, which can then be cast
// back to whatever type is appropriate by the calling code.
func List(s Store, ref KeySaver) ([]KeySaver, error) {
	keys, err := s.Keys()
	if err != nil {
		return nil, err
	}
	res := make([]KeySaver, len(keys))
	for i, k := range keys {
		v := ref.New()
		ok, err := load(s, v, k, true)
		if !ok {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

// Load fetches the backing value of k from s.  The bool indicates
// whether the value was loaded, and error contains the last error
// that occurred during the load process.
func Load(s Store, k KeySaver) (bool, error) {
	return load(s, k, k.Key(), true)
}

// Remove removes k from s.  The bool indicates whether the value was
// removed, and the error contains the last error that occurred during
// the removal process.
func Remove(s Store, k KeySaver) (bool, error) {
	if h, ok := k.(BeforeDeleteHooker); ok {
		if err := h.BeforeDelete(); err != nil {
			return false, err
		}
	}
	if err := s.Remove(k.Key()); err != nil {
		return false, err
	}
	if h, ok := k.(AfterDeleteHooker); ok {
		h.AfterDelete()
	}
	return true, nil
}

func save(s Store, k KeySaver) (bool, error) {
	if h, ok := k.(BeforeSaveHooker); ok {
		if err := h.BeforeSave(); err != nil {
			return false, err
		}
	}
	toSave := k
	if alt, ok := k.(SaveCleanHooker); ok {
		toSave = alt.SaveClean()
	}
	if err := s.Save(toSave.Key(), toSave); err != nil {
		return false, err
	}
	if h, ok := k.(AfterSaveHooker); ok {
		h.AfterSave()
	}

	return true, nil
}

// Save saves k in s, overwriting anything else that may be there.
// The bool indicates that the object was saved, and the error
// contains the last error that occurred..
func Save(s Store, k KeySaver) (bool, error) {
	return save(s, k)
}

// Create saves k in s, with the caveat that k must not already be
// present in s.  The bool indicates that the object was saved, and
// the error indicates the last error that occurred.
func Create(s Store, k KeySaver) (bool, error) {
	v := k.New()
	if ok, _ := load(s, v, k.Key(), false); ok {
		return false, fmt.Errorf("Create: thing %s:%s already exists", k.Prefix(), k.Key())
	}
	if h, ok := k.(CreateHooker); ok {
		if err := h.OnCreate(); err != nil {
			return false, err
		}
	}
	return save(s, k)
}

// Update saves k in s, with the caveat that s must already contain an
// older version of k.  If k implements ChangeHooker, then it will be
// called with the version that already exists in the backing store.
func Update(s Store, k KeySaver) (bool, error) {
	v := k.New()
	if ok, _ := load(s, v, k.Key(), false); !ok {
		return false, fmt.Errorf("Update: %s:%s does not already exist", k.Prefix(), k.Key())
	}
	if h, ok := k.(ChangeHooker); ok {
		if err := h.OnChange(v); err != nil {
			return false, err
		}
	}
	return save(s, k)
}

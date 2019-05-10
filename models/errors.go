package models

import (
	"fmt"
	"path"
	"strings"
)

// ErrorAdder is an interface that the various models that can collect
// errors for later repoting can satisfy.
type ErrorAdder interface {
	Errorf(string, ...interface{})
	AddError(error)
	HasError() error
}

// Error is the common Error type the API returns for any error
// conditions.
//
// swagger:model
type Error struct {
	Object Model `json:"-"`
	Model  string
	Key    string
	Type   string
	// Messages are any additional messages related to this Error
	Messages []string
	// code is the HTTP status code that should be used for this Error
	Code int
}

// NewError creates a new Error with a few key parameters
// pre-populated.
func NewError(t string, code int, m string) *Error {
	return &Error{Type: t, Code: code, Messages: []string{m}}
}

// Errorf appends a new error message into the Messages tracked by the
// Error.
func (e *Error) Errorf(s string, args ...interface{}) {
	if e.Messages == nil {
		e.Messages = []string{}
	}
	e.Messages = append(e.Messages, fmt.Sprintf(s, args...))
}

// Error satifies the global error interface.
func (e *Error) Error() string {
	var res string
	if e.Key != "" {
		res = fmt.Sprintf("%s: %s", e.Type, path.Join(e.Model, e.Key))
	} else if e.Model != "" {
		res = fmt.Sprintf("%s: %s", e.Type, e.Model)
	} else {
		res = fmt.Sprintf("%s", e.Type)
	}
	switch len(e.Messages) {
	case 0:
		return res
	case 1:
		return res + ": " + e.Messages[0]
	default:
		allMsgs := strings.Join(e.Messages, "\n  ")
		return res + "\n  " + allMsgs
	}
}

func (e *Error) ContainsError() bool {
	return e != nil && len(e.Messages) != 0
}

func (e *Error) AddError(src error) {
	if src == nil {
		return
	}
	if e.Messages == nil {
		e.Messages = []string{}
	}
	switch other := src.(type) {
	case *Error:
		if other.Messages != nil {
			e.Messages = append(e.Messages, other.Messages...)
		}
	case *Validation:
		if other != nil && len(other.Errors) > 0 {
			e.Messages = append(e.Messages, other.Errors...)
		}
	default:
		e.Messages = append(e.Messages, src.Error())
	}
}

func (e *Error) HasError() error {
	if e.Object != nil {
		e.Model = e.Object.Prefix()
		e.Key = e.Object.Key()
	}
	if e.ContainsError() {
		return e
	}
	return nil
}

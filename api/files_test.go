package api

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/digitalrebar/provision/models"
)

func TestFiles(t *testing.T) {
	tests := []crudTest{
		{
			name:      "list files",
			expectRes: []string{"drpcli.amd64.linux", "drpcli.amd64.windows", "drpcli.arm64.linux", "jq", "plugin_providers/"},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListBlobs("files")
			},
		},
		{
			name:      "get nonexistent file",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "files",
				Key:      "/foo",
				Type:     "GET",
				Messages: []string{"Not a regular file"},
				Code:     404,
			},
			op: func() (interface{}, error) {
				return nil, session.GetBlob(ioutil.Discard, "files", "foo")
			},
		},
		{
			name:      "create new file",
			expectRes: &models.BlobInfo{Path: "/foo", Size: 18},
			expectErr: nil,
			op: func() (interface{}, error) {
				buf := bytes.NewBufferString("Hi i am a new file")
				return session.PostBlob(buf, "files", "foo")
			},
		},
		{
			name:      "create new file in a dir (fail with name already exists)",
			expectRes: &models.BlobInfo{Path: "/foo/bar", Size: 18},
			expectErr: &models.Error{
				Model: "files",
				Key:   "/foo/bar",
				Type:  "POST",
				Messages: []string{
					"Cannot create directory /foo",
				},
				Code: 409,
			},
			op: func() (interface{}, error) {
				buf := bytes.NewBufferString("Hi i am a new file")
				return session.PostBlob(buf, "files", "foo", "bar")
			},
		},
		{
			name:      "create new file in a dir (success)",
			expectRes: &models.BlobInfo{Path: "/bar/foo", Size: 18},
			expectErr: nil,
			op: func() (interface{}, error) {
				buf := bytes.NewBufferString("Hi i am a new file")
				return session.PostBlob(buf, "files", "bar", "foo")
			},
		},
		{
			name: "list files again",
			expectRes: []string{
				"bar/",
				"drpcli.amd64.linux",
				"drpcli.amd64.windows",
				"drpcli.arm64.linux",
				"foo",
				"jq",
				"plugin_providers/",
			},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListBlobs("files")
			},
		},
		{
			name: "list files in /bar",
			expectRes: []string{
				"foo",
			},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListBlobs("files", "path", "/bar")
			},
		},
		{
			name: "delete /bar (and fail)",
			expectRes: []string{
				"foo",
			},
			expectErr: &models.Error{
				Model:    "files",
				Key:      "/bar",
				Type:     "DELETE",
				Messages: []string{"Unable to delete"},
				Code:     404,
			},
			op: func() (interface{}, error) {
				return nil, session.DeleteBlob("files", "bar")
			},
		},
		{
			name:      "delete /bar/foo",
			expectRes: nil,
			expectErr: nil,
			op: func() (interface{}, error) {
				return nil, session.DeleteBlob("files", "bar", "foo")
			},
		},
		{
			name:      "delete /bar",
			expectRes: nil,
			expectErr: nil,
			op: func() (interface{}, error) {
				return nil, session.DeleteBlob("files", "bar")
			},
		},
	}
	for _, test := range tests {
		test.run(t)
	}
}

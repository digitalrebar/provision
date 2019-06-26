package api

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/digitalrebar/provision/models"
)

func TestIsos(t *testing.T) {
	tests := []crudTest{
		{
			name:      "list isos",
			expectRes: []string{},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListBlobs("isos")
			},
		},
		{
			name:      "get nonexistent iso",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "isos",
				Key:      "/foo",
				Type:     "GET",
				Messages: []string{"Not a regular file"},
				Code:     404,
			},
			op: func() (interface{}, error) {
				return nil, session.GetBlob(ioutil.Discard, "isos", "foo")
			},
		},
		{
			name:      "create new iso",
			expectRes: &models.BlobInfo{Path: "foo", Size: 17},
			expectErr: nil,
			op: func() (interface{}, error) {
				buf := bytes.NewBufferString("Hi i am a new iso")
				return session.PostBlob(buf, "isos", "foo")
			},
		},
		{
			name: "list isos again",
			expectRes: []string{
				"foo",
			},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListBlobs("isos")
			},
		},
		{
			name:      "delete iso",
			expectRes: nil,
			expectErr: nil,
			op: func() (interface{}, error) {
				return nil, session.DeleteBlob("isos", "foo")
			},
		},
		{
			name:      "delete iso again",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "isos",
				Key:      "/foo",
				Type:     "DELETE",
				Messages: []string{"Unable to delete"},
				Code:     404,
			},
			op: func() (interface{}, error) {
				return nil, session.DeleteBlob("isos", "foo")
			},
		},
	}
	for _, test := range tests {
		test.run(t)
	}
}

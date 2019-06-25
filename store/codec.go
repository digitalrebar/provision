package store

import (
	"encoding/json"

	"github.com/ghodss/yaml"
)

type codec struct {
	enc func(interface{}) ([]byte, error)
	dec func([]byte, interface{}) error
	ext string
}

func (c *codec) Encode(i interface{}) ([]byte, error) {
	return c.enc(i)
}

func (c *codec) Decode(buf []byte, i interface{}) error {
	return c.dec(buf, i)
}

func (c *codec) Ext() string {
	return c.ext
}

// Codec is responsible for encoding and decoding raw Go objects into
// a serializable form.
type Codec interface {
	// Encode takes an object and turns it into a byte array.
	Encode(interface{}) ([]byte, error)
	// Decode takes a byte array and decodes it into an object
	Decode([]byte, interface{}) error
	// Ext is the file extension that should be used for this encoding
	// if we are encoding to a filesystem.
	Ext() string
}

// JsonCodec implements Codec for encoding/decoding to JSON.
var JsonCodec = &codec{
	enc: json.Marshal,
	dec: json.Unmarshal,
	ext: ".json",
}

func yamlDecode(buf []byte, d interface{}) error {
	return yaml.Unmarshal(buf, d)
}

// YamlCodec implements a Codec for encoding/decoding to YAML
var YamlCodec = &codec{
	enc: yaml.Marshal,
	dec: yamlDecode,
	ext: ".yaml",
}

var DefaultCodec = JsonCodec

package converter

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Marshaler interface {
	Tag() string
	Ext() string
	Marshal(v interface{}) ([]byte, error)
}
type YamlMarshaler struct {
}

func (YamlMarshaler) Tag() string {
	return `yaml`
}
func (YamlMarshaler) Ext() string {
	return `.yaml`
}
func (YamlMarshaler) Marshal(v interface{}) (out []byte, err error) {
	return yaml.Marshal(v)
}

type JsonMarshaler struct {
}

func (JsonMarshaler) Tag() string {
	return `json`
}
func (JsonMarshaler) Ext() string {
	return `.json`
}
func (JsonMarshaler) Marshal(v interface{}) (out []byte, err error) {
	return json.Marshal(v)
}

type PrettyJsonMarshaler struct {
	JsonMarshaler
}

func (PrettyJsonMarshaler) Marshal(v interface{}) (out []byte, err error) {
	return json.MarshalIndent(v, "", "\t")
}

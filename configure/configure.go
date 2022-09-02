package configure

import (
	"encoding/json"
	"path/filepath"

	"github.com/google/go-jsonnet"
)

type Configure struct {
	Version   string     `json:"version"`
	Endpoints []Endpoint `json:"endpoints"`
}

func (c *Configure) String() string {
	if c == nil {
		return "null"
	}
	b, e := json.MarshalIndent(c, ``, `	`)
	if e != nil {
		return e.Error()
	}
	return string(b)
}

func (c *Configure) Load(filename string) (e error) {
	vm := jsonnet.MakeVM()
	jsonStr, e := vm.EvaluateFile(filename)
	if e != nil {
		return
	}
	e = json.Unmarshal([]byte(jsonStr), c)
	if e != nil {
		return
	}
	return
}

type Endpoint struct {
	Output    string   `json:"output"`
	Target    string   `json:"target"`
	Source    string   `json:"source"`
	Resources []string `json:"resources"`
	Prefix    string   `json:"-"`
}

func (ep *Endpoint) Format(dir string) (e error) {
	if len(ep.Resources) == 0 {
		return
	}
	ep.Output, e = ep.abs(dir, ep.Output)
	if e != nil {
		return
	}
	ep.Target, e = ep.abs(dir, ep.Target)
	if e != nil {
		return
	}
	ep.Source, e = ep.abs(dir, ep.Source)
	if e != nil {
		return
	}
	ep.Prefix = ep.Source + string(filepath.Separator)
	return
}
func (ep *Endpoint) abs(dir, src string) (dst string, e error) {
	if src == `` {
		dst = dir
		return
	}
	if filepath.IsAbs(src) {
		dst = filepath.Clean(src)
	} else {
		dst = filepath.Clean(filepath.Join(dir, src))
	}
	return
}

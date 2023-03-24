package configure

import (
	"encoding/json"
	"path/filepath"
	"runtime"

	"github.com/google/go-jsonnet"
	"github.com/powerpuffpenguin/ejt/internal/fix"
	"github.com/powerpuffpenguin/ejt/version"
)

type Configure struct {
	Version   string     `json:"version"`
	Endpoints []Endpoint `json:"endpoints"`
	ExtStrs   []string   `json:"ext_strs"`
	JPath     []string   `json:"jpath"`
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

func (c *Configure) Load(dir, filename string, jpath0 []string) (e error) {
	vm := jsonnet.MakeVM()
	vm.Importer(&fix.FileImporter{})
	jsonStr, e := vm.EvaluateFile(filepath.Join(dir, filename))
	if e != nil {
		return
	}
	e = json.Unmarshal([]byte(jsonStr), c)
	if e != nil {
		return
	}
	var (
		keys  = make(map[string]bool)
		jpath []string
	)
	for _, s := range jpath0 {
		if filepath.IsAbs(s) {
			s = filepath.Clean(s)
		} else {
			s = filepath.Join(dir, s)
		}
		if keys[s] {
			continue
		}
		keys[s] = true
		jpath = append(jpath, s)
	}
	for _, s := range c.JPath {
		if filepath.IsAbs(s) {
			s = filepath.Clean(s)
		} else {
			s = filepath.Join(dir, s)
		}
		if keys[s] {
			continue
		}
		keys[s] = true
		jpath = append(jpath, s)
	}
	c.JPath = jpath
	c.ExtStrs = append([]string{
		`dev=0`,
		`ejt.version=` + version.Version,
		`ejt.os=` + runtime.GOOS,
		`ejt.arch=` + runtime.GOARCH,
		`ejt.go_version=` + runtime.Version(),
		`ejt.dir=` + dir,
		`ejt.jsonnet=` + jsonnet.Version(),
	}, c.ExtStrs...)
	for i := 0; i < len(c.Endpoints); i++ {
		endpoint := &c.Endpoints[i]
		e = c.format(dir, endpoint, jpath)
		if e != nil {
			return
		}
	}
	return
}
func (c *Configure) format(dir string, endpoint *Endpoint, jpath0 []string) (e error) {
	if len(endpoint.Resources) == 0 {
		return
	}
	jpath := append([]string{}, jpath0...)

	keys := make(map[string]bool, len(jpath)+len(endpoint.JPath))
	for _, k := range jpath {
		keys[k] = true
	}

	for _, s := range endpoint.JPath {
		if filepath.IsAbs(s) {
			s = filepath.Clean(s)
		} else {
			s = filepath.Join(dir, s)
		}
		if keys[s] {
			continue
		}
		keys[s] = true
		jpath = append(jpath, s)
	}
	endpoint.JPath = jpath
	return
}

type Endpoint struct {
	Output    string   `json:"output"`
	Target    string   `json:"target"`
	Source    string   `json:"source"`
	Resources []string `json:"resources"`
	Prefix    string   `json:"-"`

	ExtStrs []string `json:"ext_strs"`
	JPath   []string `json:"jpath"`
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

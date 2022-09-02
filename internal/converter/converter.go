package converter

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/powerpuffpenguin/ejt/configure"
	"gopkg.in/yaml.v3"
)

type Converter struct {
	cnf *configure.Configure
}

func New() (c *Converter, e error) {
	pwd, e := filepath.Abs(".")
	if e != nil {
		return
	}
	var (
		cnf   configure.Configure
		info  fs.FileInfo
		found bool
	)
	for {
		path := filepath.Join(pwd, `ejt.jsonnet`)
		info, e = os.Stat(path)
		if e != nil {
			if os.IsNotExist(e) {
				dir := filepath.Dir(pwd)
				if dir == pwd {
					break
				}
				pwd = dir
				continue
			}
			return
		} else if info.IsDir() {
			dir := filepath.Dir(pwd)
			if dir == pwd {
				break
			}
			pwd = dir
			continue
		}

		e = cnf.Load(filepath.Join(pwd, `ejt.jsonnet`))
		if e != nil {
			return
		}
		found = true
		break
	}
	if !found {
		e = errors.New(`not found ejt.jsonnet`)
		return
	}
	for i := 0; i < len(cnf.Endpoints); i++ {
		e = cnf.Endpoints[i].Format(pwd)
		if e != nil {
			return
		}
	}

	c = &Converter{
		cnf: &cnf,
	}
	return
}
func (c *Converter) Yaml(test, move, copy bool) {
	cnf := c.cnf
	for i := 0; i < len(cnf.Endpoints); i++ {
		endpoint := &cnf.Endpoints[i]
		if len(endpoint.Resources) != 0 {
			log.Printf("yaml endpoints[%d] from %s\n", i, endpoint.Source)
			for j, resource := range endpoint.Resources {
				log.Printf(" - %-2d %s\n", j, resource)
				c.yaml(endpoint, resource, test)
			}

			if move {
				for j, resource := range endpoint.Resources {
					log.Printf(" - %-2d move %s\n", j, resource)
					c.move(endpoint, resource, test)
				}
			} else if copy {
				for j, resource := range endpoint.Resources {
					log.Printf(" - %-2d copy %s\n", j, resource)
					c.copy(endpoint, resource, test)
				}
			}
		}
	}
}
func (c *Converter) yaml(endpoint *configure.Endpoint, resource string, test bool) {
	filename := filepath.Clean(filepath.Join(endpoint.Source, resource))
	if !strings.HasPrefix(filename, endpoint.Prefix) {
		log.Fatalln("resource illegal")
	}
	vm := jsonnet.MakeVM()
	jsonStr, e := vm.EvaluateFile(filename)
	if e != nil {
		log.Fatalln(e)
	}
	var m interface{}
	e = json.Unmarshal([]byte(jsonStr), &m)
	if e != nil {
		log.Fatalln(e)
	}
	b, e := yaml.Marshal(&m)
	if e != nil {
		log.Fatalln(e)
	}
	if test {
		fmt.Println(string(b))
		return
	}

	name := filename[len(endpoint.Prefix):]
	ext := filepath.Ext(name)
	if ext != `` {
		name = name[:len(name)-len(ext)]
	}
	name += `.yaml`
	output := filepath.Join(endpoint.Output, name)
	os.MkdirAll(filepath.Dir(output), 0775)
	e = ioutil.WriteFile(output, b, 0644)
	if e != nil {
		log.Fatalln(e)
	}
}
func (c *Converter) get(endpoint *configure.Endpoint, resource string) (src, dst string) {
	filename := filepath.Clean(filepath.Join(endpoint.Source, resource))
	if !strings.HasPrefix(filename, endpoint.Prefix) {
		log.Fatalln("resource illegal")
	}
	name := filename[len(endpoint.Prefix):]
	ext := filepath.Ext(name)
	if ext != `` {
		name = name[:len(name)-len(ext)]
	}
	name += `.yaml`

	src = filepath.Join(endpoint.Output, name)
	dst = filepath.Join(endpoint.Target, name)
	return
}
func (c *Converter) md5(filename string) (b []byte, e error) {
	f, e := os.Open(filename)
	if e != nil {
		return
	}
	h := md5.New()
	_, e = io.Copy(h, f)
	f.Close()
	if e != nil {
		return
	}
	b = h.Sum(nil)
	return
}
func (c *Converter) compare(src, dst string) (changed bool) {
	l, e := c.md5(src)
	if e != nil {
		if os.IsNotExist(e) {
			changed = true
			return
		}
		log.Fatalln(e)
	}
	r, e := c.md5(dst)
	if e != nil {
		log.Fatalln(e)
	}
	changed = !bytes.Equal(l, r)
	return
}
func (c *Converter) move(endpoint *configure.Endpoint, resource string, test bool) {
	src, dst := c.get(endpoint, resource)
	if test {
		log.Println(` # move`, src, `->`, dst)
	} else {
		if !c.compare(src, dst) {
			log.Println(` # not changed`, resource)
			return
		}

		os.MkdirAll(filepath.Dir(dst), 0775)
		e := os.Rename(src, dst)
		if e != nil {
			log.Fatalln(e)
		}
	}
}
func (c *Converter) copy(endpoint *configure.Endpoint, resource string, test bool) {
	src, dst := c.get(endpoint, resource)
	if test {
		log.Println(` # copy`, src, `->`, dst)
	} else {
		if !c.compare(src, dst) {
			log.Println(` # not changed`, resource)
			return
		}

		r, e := os.OpenFile(src, os.O_RDONLY, 0664)
		if e != nil {
			log.Fatalln(e)
		}
		defer r.Close()
		os.MkdirAll(filepath.Dir(dst), 0775)
		w, e := os.Create(dst)
		if e != nil {
			log.Fatalln(e)
		}
		defer w.Close()
		_, e = io.Copy(w, r)
		if e != nil {
			log.Fatalln(e)
		}
	}
}

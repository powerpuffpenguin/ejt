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
	"github.com/powerpuffpenguin/ejt/internal/fix"
)

func parseExtVar(str string, exists map[string]string) (k, v string) {
	index := strings.Index(str, `=`)
	if index == -1 {
		k = str
		var found bool
		if v, found = exists[k]; found {
			return
		}
		v = os.Getenv(str)
	} else {
		k = str[:index]
		v = str[index+1:]
	}
	return
}

type Converter struct {
	cnf     *configure.Configure
	keyStrs map[string]string
}

func New(extStrs, jpath []string) (c *Converter, e error) {
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

		e = cnf.Load(pwd, `ejt.jsonnet`, jpath)
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

	keyStrs := make(map[string]string, len(cnf.ExtStrs)+len(extStrs))
	for _, str := range cnf.ExtStrs {
		k, v := parseExtVar(str, nil)
		keyStrs[k] = v
	}
	for _, str := range extStrs {
		k, v := parseExtVar(str, nil)
		keyStrs[k] = v
	}
	c = &Converter{
		cnf:     &cnf,
		keyStrs: keyStrs,
	}
	return
}

func (c *Converter) vm(endpoint *configure.Endpoint) (vm *jsonnet.VM) {
	vm = jsonnet.MakeVM()
	vm.Importer(&fix.FileImporter{
		JPaths: endpoint.JPath,
	})
	for k, v := range c.keyStrs {
		vm.ExtVar(k, v)
	}
	for _, str := range endpoint.ExtStrs {
		vm.ExtVar(parseExtVar(str, nil))
	}
	return
}
func (c *Converter) Convert(marshaler Marshaler, test, move, copy, replace bool) {
	var (
		cnf   = c.cnf
		hashs [][]byte
		hash  []byte
	)
	for i := 0; i < len(cnf.Endpoints); i++ {
		endpoint := &cnf.Endpoints[i]
		if len(endpoint.Resources) != 0 {
			log.Printf("%s endpoints[%d] from %s\n",
				marshaler.Tag(),
				i, endpoint.Source,
			)
			if test {
				hashs = make([][]byte, len(endpoint.Resources))
			}
			vm := c.vm(endpoint)
			for j, resource := range endpoint.Resources {
				log.Printf(" - %-2d %s\n", j, resource)
				hash = c.convert(vm, marshaler, endpoint, resource, test)
				if test {
					hashs[j] = hash
				}
			}

			if move {
				for j, resource := range endpoint.Resources {
					log.Printf(" - %-2d move %s\n", j, resource)
					if test {
						hash = hashs[j]
					}
					c.move(marshaler, endpoint, resource, hash, test, replace)
				}
			} else if copy {
				for j, resource := range endpoint.Resources {
					log.Printf(" - %-2d copy %s\n", j, resource)
					if test {
						hash = hashs[j]
					}
					c.copy(marshaler, endpoint, resource, hash, test, replace)
				}
			}
		}
	}
}

func (c *Converter) convert(vm *jsonnet.VM, marshaler Marshaler, endpoint *configure.Endpoint, resource string, test bool) (hash []byte) {
	filename := filepath.Clean(filepath.Join(endpoint.Source, resource))
	if !strings.HasPrefix(filename, endpoint.Prefix) {
		log.Fatalln("resource illegal")
	}

	jsonStr, e := vm.EvaluateFile(filename)
	if e != nil {
		log.Fatalln(e)
	}
	var m interface{}
	e = json.Unmarshal([]byte(jsonStr), &m)
	if e != nil {
		log.Fatalln(e)
	}
	b, e := marshaler.Marshal(&m)
	if e != nil {
		log.Fatalln(e)
	}
	if test {
		fmt.Println(string(b))
		h := md5.Sum(b)
		hash = h[:]
		return
	}

	name := filename[len(endpoint.Prefix):]
	ext := filepath.Ext(name)
	if ext != `` {
		name = name[:len(name)-len(ext)]
	}
	name += marshaler.Ext()
	output := filepath.Join(endpoint.Output, name)
	os.MkdirAll(filepath.Dir(output), 0775)
	e = ioutil.WriteFile(output, b, 0644)
	if e != nil {
		log.Fatalln(e)
	}
	return
}
func (c *Converter) get(marshaler Marshaler, endpoint *configure.Endpoint, resource string) (src, dst string) {
	filename := filepath.Clean(filepath.Join(endpoint.Source, resource))
	if !strings.HasPrefix(filename, endpoint.Prefix) {
		log.Fatalln("resource illegal")
	}
	name := filename[len(endpoint.Prefix):]
	ext := filepath.Ext(name)
	if ext != `` {
		name = name[:len(name)-len(ext)]
	}
	name += marshaler.Ext()

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
	l, e := c.md5(dst)
	if e != nil {
		if os.IsNotExist(e) {
			changed = true
			return
		}
		log.Fatalln(e)
	}
	r, e := c.md5(src)
	if e != nil {
		log.Fatalln(e)
	}
	changed = !bytes.Equal(l, r)
	return
}
func (c *Converter) compareHash(src []byte, dst string) (changed bool) {
	l, e := c.md5(dst)
	if e != nil {
		if os.IsNotExist(e) {
			changed = true
			return
		}
		log.Fatalln(e)
	}
	changed = !bytes.Equal(l, src)
	return
}
func (c *Converter) move(marshaler Marshaler, endpoint *configure.Endpoint, resource string, hash []byte, test, replace bool) {
	src, dst := c.get(marshaler, endpoint, resource)
	if test {
		if !replace && !c.compareHash(hash, dst) {
			log.Println(` # not changed`, resource)
			return
		}
	} else {
		if !replace && !c.compare(src, dst) {
			log.Println(` # not changed`, resource)
			return
		}
		r, e := os.OpenFile(src, os.O_RDONLY, 0664)
		if e != nil {
			log.Fatalln(e)
		}
		defer r.Close()
		os.MkdirAll(filepath.Dir(dst), 0775)
		temp := dst + `.temp`
		w, e := os.Create(temp)
		if e != nil {
			log.Fatalln(e)
		}
		_, e = io.Copy(w, r)
		w.Close()
		if e != nil {
			log.Fatalln(e)
		}

		e = os.Rename(temp, dst)
		if e != nil {
			log.Fatalln(e)
		}
	}
}
func (c *Converter) copy(marshaler Marshaler, endpoint *configure.Endpoint, resource string, hash []byte, test, replace bool) {
	src, dst := c.get(marshaler, endpoint, resource)
	if test {
		if !replace && !c.compareHash(hash, dst) {
			log.Println(` # not changed`, resource)
			return
		}
	} else {
		if !replace && !c.compare(src, dst) {
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

package converter

import (
	"os"
	"path/filepath"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

var nativeFuncs = []*jsonnet.NativeFunction{
	{
		Name: "os.readText",
		Params: ast.Identifiers{
			"filename",
		},
		Func: func(i []interface{}) (interface{}, error) {
			b, e := os.ReadFile(i[0].(string))
			if e != nil {
				return nil, e
			}
			return string(b), nil
		},
	},
	{
		Name: "filepath.join",
		Params: ast.Identifiers{
			"dir", "name",
		},
		Func: func(i []interface{}) (interface{}, error) {
			dir := i[0].(string)
			name := i[1].(string)
			return filepath.Join(dir, name), nil
		},
	},
	{
		Name: "filepath.clean",
		Params: ast.Identifiers{
			"path",
		},
		Func: func(i []interface{}) (interface{}, error) {
			return filepath.Clean(i[0].(string)), nil
		},
	},
	{
		Name: "filepath.abs",
		Params: ast.Identifiers{
			"path",
		},
		Func: func(i []interface{}) (interface{}, error) {
			return filepath.Abs(i[0].(string))
		},
	},
	{
		Name: "filepath.isAbs",
		Params: ast.Identifiers{
			"path",
		},
		Func: func(i []interface{}) (interface{}, error) {
			return filepath.IsAbs(i[0].(string)), nil
		},
	},
	{
		Name: "filepath.base",
		Params: ast.Identifiers{
			"path",
		},
		Func: func(i []interface{}) (interface{}, error) {
			return filepath.Base(i[0].(string)), nil
		},
	},
	{
		Name: "filepath.dir",
		Params: ast.Identifiers{
			"path",
		},
		Func: func(i []interface{}) (interface{}, error) {
			return filepath.Dir(i[0].(string)), nil
		},
	},
	{
		Name: "filepath.ext",
		Params: ast.Identifiers{
			"path",
		},
		Func: func(i []interface{}) (interface{}, error) {
			return filepath.Ext(i[0].(string)), nil
		},
	},
}

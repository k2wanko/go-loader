package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gopherjs/gopherjs/js"
)

var (
	exec        *js.Object
	loaderUtils *js.Object
)

func init() {
	exec = js.Global.Call("require", "child_process").Get("execSync")
	loaderUtils = js.Global.Call("require", "loader-utils")
}

func compile(src []byte) (code io.Reader, srcMap io.Reader, closer func(), err error) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	err = ioutil.WriteFile(filepath.Join(dir, "main.go"), src, 0777)
	if err != nil {
		return
	}

	prevDir, _ := filepath.Abs(".")
	os.Chdir(dir)
	defer os.Chdir(prevDir)

	name := filepath.Join(dir, "a.js")
	exec.Invoke("gopherjs build -o a.js")

	s, err := os.Open(name)
	if err != nil {
		return
	}
	m, err := os.Open(name + ".map")
	if err != nil {
		s.Close()
		return
	}

	return s, m, func() {
		s.Close()
		m.Close()
	}, nil
}

func main() {
	js.Module.Set("exports", js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
		this.Call("cacheable")
		cb := this.Call("async")
		go func() {
			src, srcMap, close, err := compile([]byte(args[0].String()))
			if err != nil {
				cb.Invoke(err)
				return
			}
			defer close()

			d, err := ioutil.ReadAll(src)
			if err != nil {
				cb.Invoke(err)
				return
			}
			m, err := ioutil.ReadAll(srcMap)
			if err != nil {
				cb.Invoke(err)
				return
			}

			cb.Invoke(nil, string(d), string(m))
		}()

		return js.Undefined
	}))
}

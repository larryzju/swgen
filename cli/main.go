package main

import (
	"flag"
	"log"
	"os"
	"path"
	"runtime/debug"

	"github.com/larryzju/swgen"
)

var (
	input    = flag.String("input", "", "input directory")
	output   = flag.String("output", "", "output directory")
	root     = flag.String("root", "/", "root directory")
	template = flag.String("template", "", "html template")
	verbose  = flag.Bool("verbose", false, "verbose")
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("stacktrace from panic: ", string(debug.Stack()))
		}
	}()

	flag.Parse()

	// ignore
	var ignore swgen.Ignore
	f, err := os.Open(path.Join(*input, ".swignore"))
	if err != nil {
		ignore = &swgen.BasicIgnore{}
	}
	defer f.Close()
	ignore = swgen.NewBasicIgnore(f)

	err = swgen.Generate(*input, *output, *root, *template, ignore)
	if err != nil {
		log.Panic(err)
	}
}

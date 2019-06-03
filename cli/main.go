package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/larryzju/swgen"
)

var (
	input   = flag.String("input", "", "input directory")
	output  = flag.String("output", "", "output directory")
	root    = flag.String("root", "/", "root directory")
	verbose = flag.Bool("verbose", false, "verbose")
)

func main() {
	flag.Parse()

	// ignore
	var ignore swgen.Ignore
	f, err := os.Open(path.Join(*input, ".swignore"))
	if err != nil {
		ignore = &swgen.BasicIgnore{}
	}
	defer f.Close()
	ignore = swgen.NewBasicIgnore(f)

	// build
	c, err := swgen.NewDirNode(*input, *input, ignore)
	if err != nil {
		log.Panic(err)
	}

	dir := c.(*swgen.DirNode)

	// generate navigator
	nav, err := dir.Navbar(*root)
	if err != nil {
		log.Fatal(err)
	}

	// generate pages
	err = dir.Generate(*output, nav)
	if err != nil {
		log.Fatal(err)
	}
}

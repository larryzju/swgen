package main

import (
	"flag"
	"log"

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
	c, err := swgen.NewDirNode(*input, *input)
	if err != nil {
		log.Fatal(err)
	}

	dir := c.(*swgen.DirNode)
	nav, err := dir.Navbar(*root)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(nav)

	err = dir.Generate(*output, nav)
	if err != nil {
		log.Fatal(err)
	}
}

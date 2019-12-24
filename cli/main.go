package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/larryzju/swgen"
)

var (
	input   = flag.String("input", ".", "input directory")
	output  = flag.String("output", "./output", "output directory")
	root    = flag.String("root", "", "root directory")
	force   = flag.Bool("force", false, "ignore timestamp")
	verbose = flag.Bool("verbose", false, "verbose")
)

func main() {
	flag.Parse()

	var ignore swgen.Ignore
	f, err := os.Open(filepath.Join(*input, ".swignore"))
	if err != nil {
		ignore = &swgen.BasicIgnore{}
	}
	defer f.Close()
	ignore = swgen.NewBasicIgnore(f)

	templatePattern := filepath.Join(*input, ".template/*.html")
	tmpl := template.Must(template.ParseGlob(templatePattern))
	log.Printf("template=%v", tmpl)
	sw := swgen.Swgen{
		URLRoot:  *root,
		Source:   *input,
		Target:   *output,
		Force:    *force,
		Ignore:   ignore,
		Template: tmpl,
	}

	err = sw.Run()
	if err != nil {
		log.Panic(err)
	}
}

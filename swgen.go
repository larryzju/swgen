package swgen

import (
	"encoding/json"
	"html/template"
	"os"
	"path"
	"time"
)

type Target struct {
	output       string
	root         string
	htmlTemplate *template.Template
}

func NewTarget(output, root, templatePath string) (target *Target, err error) {
	tmpl, err := template.New("page.html").ParseFiles(templatePath)
	if err != nil {
		return
	}
	return &Target{output, root, tmpl}, nil
}

func (t *Target) Root() string                     { return t.output }
func (t *Target) URLRoot() string                  { return t.root }
func (t *Target) HtmlTemplate() *template.Template { return t.htmlTemplate }

type Metadata struct {
	navigator  string
	buildTime  time.Time
	gitVersion string
}

func Generate(input, output, root, templatePath string, ignore Ignore) (err error) {
	source := &Source{input, input, ignore}
	dir, err := NewDir(source)
	if err != nil {
		return
	}

	target, err := NewTarget(output, root, templatePath)
	if err != nil {
		return
	}

	navigator, err := dir.generateNavigator(target)
	if err != nil {
		return
	}

	bytes, err := json.Marshal(navigator)
	if err != nil {
		return
	}

	navPath := path.Join(output, "nav.json")
	nav, err := os.Create(navPath)
	if err != nil {
		return
	}

	_, err = nav.Write(bytes)
	if err != nil {
		return
	}

	metadata := &Metadata{navPath, time.Now(), ""}
	err = dir.Flush(metadata, target)
	if err != nil {
		return
	}

	return
}

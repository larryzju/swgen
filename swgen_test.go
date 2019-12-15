package swgen

import (
	"html/template"
	"strings"
	"testing"
)

type dummyIgnore struct{}

func (i *dummyIgnore) Ignore(path string) bool {
	if strings.HasPrefix(path, ".git") {
		return true
	}

	if strings.HasPrefix(path, "testing") {
		return true
	}

	if strings.HasSuffix(path, "~") {
		return true
	}
	
	return false
}

func TestGenerate(t *testing.T) {
	tmpl := template.Must(template.New("page").Parse(`Page: {{.Page}}`))
	sw := Swgen{
		source:   ".",
		target:   "testing",
		ignore:   &dummyIgnore{},
		template: tmpl,
	}
	if err := sw.Run(); err != nil {
		t.Error(err)
	}
}

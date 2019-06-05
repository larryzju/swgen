package node

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

type Tmpl struct {
	rel  string
	path string
	info os.FileInfo
}

type TmplData interface {
	Title() string
	URLRoot() string
	Navigator() template.HTML
	Body() template.HTML
}

type TmplDataWrapper struct {
	Node
	Metadata
	Target
}

func wrapTmplData(n Node, m Metadata, t Target) *TmplDataWrapper {
	return &TmplDataWrapper{Node: n, Target: t, Metadata: m}
}

func NewTmpl(s Source) (page Node, err error) {
	root := s.Root()
	path := s.Path()
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return
	}

	return &Tmpl{
		rel:  rel,
		path: path,
		info: info,
	}, nil
}

func (p *Tmpl) Title() string {
	return p.info.Name()
}

func (p *Tmpl) Rel() string {
	re := regexp.MustCompile(".tmpl$")
	return re.ReplaceAllString(p.rel, "")
}

func (p *Tmpl) LastUpdate() time.Time {
	return p.info.ModTime()
}

func (p *Tmpl) String() string {
	return p.Rel()
}

func (p *Tmpl) Flush(m Metadata, t Target) (err error) {
	tmpl, err := template.ParseFiles(p.path)
	if err != nil {
		return err
	}

	dstpath := path.Join(t.Root(), p.Rel())
	dst, err := os.Create(dstpath)
	if err != nil {
		return err
	}
	defer dst.Close()

	data := wrapTmplData(p, m, t)
	return tmpl.Execute(dst, data)
}

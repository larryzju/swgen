package node

import (
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

type HtmlNode struct {
	rel  string
	path string
	info os.FileInfo
}

type HtmlDataWrapper struct {
	TmplDataWrapper
	Body template.HTML
}

func wrapHtmlData(n Node, m Metadata, t Target, b template.HTML) *HtmlDataWrapper {
	return &HtmlDataWrapper{TmplDataWrapper{Node: n, Target: t, Metadata: m}, b}
}

func NewHtml(s Source) (node Node, err error) {
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

	return &HtmlNode{
		rel:  rel,
		path: path,
		info: info,
	}, nil
}

func (p *HtmlNode) Title() string {
	return p.info.Name()
}

func (p *HtmlNode) Rel() string {
	re := regexp.MustCompile(".org$")
	return re.ReplaceAllString(p.rel, ".html")
}

func (p *HtmlNode) LastUpdate() time.Time {
	return p.info.ModTime()
}

func (p *HtmlNode) String() string {
	return p.Rel()
}

func (p *HtmlNode) Flush(m Metadata, t Target) (err error) {
	file, err := os.Open(p.path)
	if err != nil {
		return
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	f, err := os.Create(path.Join(t.Root(), p.Rel()))
	if err != nil {
		return
	}

	data := wrapHtmlData(p, m, t, template.HTML(bytes))
	tmpl := t.HtmlTemplate()
	return tmpl.Execute(f, data)
}

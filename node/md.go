package node

import (
	"html/template"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

type MarkdownNode struct {
	rel  string
	path string
	info os.FileInfo
}

func NewMarkdown(s Source) (node Node, err error) {
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

	return &MarkdownNode{
		rel:  rel,
		path: path,
		info: info,
	}, nil
}

func (p *MarkdownNode) Title() string {
	return p.info.Name()
}

func (p *MarkdownNode) Rel() string {
	re := regexp.MustCompile(".md$")
	return re.ReplaceAllString(p.rel, ".html")
}

func (p *MarkdownNode) LastUpdate() time.Time {
	return p.info.ModTime()
}

func (p *MarkdownNode) String() string {
	return p.Rel()
}

func (p *MarkdownNode) Flush(m Metadata, t Target) (err error) {
	body, err := p.render()
	if err != nil {
		return
	}

	f, err := os.Create(path.Join(t.Root(), p.Rel()))
	if err != nil {
		return
	}

	data := wrapHTMLData(p, m, t, body)
	return htmlTemplate.Execute(f, data)
}

func (p *MarkdownNode) render() (c template.HTML, err error) {
	cmd := exec.Command("kramdown", p.path)
	bytes, err := cmd.Output()
	if err != nil {
		return
	}

	return template.HTML(bytes), nil
}

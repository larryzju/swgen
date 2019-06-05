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

type OrgNode struct {
	rel  string
	path string
	info os.FileInfo
}

func NewOrg(s Source) (node Node, err error) {
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

	return &OrgNode{
		rel:  rel,
		path: path,
		info: info,
	}, nil
}

func (p *OrgNode) Title() string {
	return p.info.Name()
}

func (p *OrgNode) Rel() string {
	re := regexp.MustCompile(".org$")
	return re.ReplaceAllString(p.rel, ".html")
}

func (p *OrgNode) LastUpdate() time.Time {
	return p.info.ModTime()
}

func (p *OrgNode) String() string {
	return p.Rel()
}

func (p *OrgNode) Flush(m Metadata, t Target) (err error) {
	body, err := p.render()
	if err != nil {
		return
	}

	f, err := os.Create(path.Join(t.Root(), p.Rel()))
	if err != nil {
		return
	}

	data := wrapHtmlData(p, m, t, body)
	tmpl := t.HtmlTemplate()
	return tmpl.Execute(f, data)
}

func (p *OrgNode) render() (c template.HTML, err error) {
	cmd := exec.Command("pandoc", "-f", "org", "-t", "html", "-i", p.path)
	bytes, err := cmd.Output()
	if err != nil {
		return
	}

	return template.HTML(bytes), nil
}

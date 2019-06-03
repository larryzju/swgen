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

func NewOrgNode(root string, src string) (node Node, err error) {
	path, err := filepath.Abs(src)
	if err != nil {
		return
	}

	info, err := os.Stat(src)
	if err != nil {
		return
	}

	rel, err := filepath.Rel(root, src)
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

func (p *OrgNode) Content() (c template.HTML, err error) {
	cmd := exec.Command("pandoc", "-f", "org", "-t", "html", "-i", p.path)
	bytes, err := cmd.Output()
	if err != nil {
		return
	}

	c = template.HTML(bytes)
	return

}

func (p *OrgNode) Link(root string) string {
	return path.Join(root, p.Rel())
}

package swgen

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var NotRenderableFile = errors.New("file can not be rendered")

var RenderFns = map[string]RenderFn{
	".md":   RenderMarkdown,
	".org":  RenderOrg,
	".html": RenderHTML,
	".htm":  RenderHTML,
}

type RenderFn func(*Node, *Metadata) (template.HTML, error)

type Node struct {
	*Swgen
	Info     os.FileInfo
	Path     string
	Children []*Node
	Home     *Node
	Next     *Node
	Prev     *Node
	Up       *Node
}

type Metadata struct {
	Root      string
	BuildTime time.Time
	Version   string
}

func (n *Node) String() string {
	return n.string(0)
}

func (n *Node) string(indent int) string {
	sb := &strings.Builder{}
	prefix := strings.Repeat("  ", indent)
	fmt.Fprintf(sb, "%s%s\n", prefix, n.Info.Name())
	for _, c := range n.Children {
		fmt.Fprintf(sb, "%s%s", prefix, c.string(indent+1))
	}
	return sb.String()
}

func (n *Node) Render(meta *Metadata) (template.HTML, error) {
	ext := filepath.Ext(n.Path)
	render, ok := RenderFns[ext]
	if !ok {
		return template.HTML(""), NotRenderableFile
	}
	return render(n, meta)
}

func RenderMarkdown(n *Node, m *Metadata) (template.HTML, error) {
	cmd := exec.Command("kramdown", n.Path)
	bytes, err := cmd.Output()
	if err != nil {
		return template.HTML(""), err
	}
	return template.HTML(bytes), nil
}

func RenderOrg(n *Node, m *Metadata) (template.HTML, error) {
	cmd := exec.Command("pandoc", "--fail-if-warnings", "-f", "org", "-t", "html", "--mathjax", "-i", filepath.Base(n.Path))
	cmd.Dir = filepath.Dir(n.Path)
	bytes, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("run command on %s failed: %s", n.Path, string(err.(*exec.ExitError).Stderr))
		return template.HTML(""), nil
	}

	return template.HTML(bytes), nil
}

func RenderHTML(n *Node, m *Metadata) (template.HTML, error) {
	bytes, err := ioutil.ReadFile(n.Path)
	if err != nil {
		return template.HTML(""), nil
	}

	return template.HTML(bytes), nil
}

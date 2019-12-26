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

var dirTemplate = template.Must(template.New("dir").Parse(`
<div>
  {{.Rel}}
  <ul>
    {{range .Children}}
    <li>
      <a href="{{.PageURL}}">{{.Name}}</a>
    </li>
    {{end}}
  </ul>
</div>
`))

var NotRenderableFile = errors.New("file can not be rendered")

var RenderFns = map[string]RenderFn{
	".md":   RenderMarkdown,
	".org":  RenderOrg,
	".html": RenderHTML,
	".htm":  RenderHTML,
}

type RenderFn func(*Node, *Metadata) (template.HTML, error)

// Node is a file or directory
type Node struct {
	*Swgen
	Info     os.FileInfo
	path     string
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

// PageURL is used to generate the HTML path
func (n *Node) PageURL() (string, error) {
	rel, err := filepath.Rel(n.Source, n.path)
	if err != nil {
		return "", err
	}

	url := filepath.Join("/", n.URLRoot, rel)
	if !n.Info.IsDir() {
		ext := filepath.Ext(n.Info.Name())
		if _, ok := RenderFns[ext]; ok {
			url += ".html"
		}
	}
	return url, nil
}

// Rel get the node's relative path
func (n *Node) Rel() (string, error) {
	return filepath.Rel(n.Source, n.path)
}

// Name get the node's relative name
func (n *Node) Name() (string, error) {
	path, err := filepath.Rel(n.Source, n.path)
	if err != nil {
		return "", err
	}

	return filepath.Base(path), nil
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

// Render page content
func (n *Node) Render(meta *Metadata) (template.HTML, error) {
	ext := filepath.Ext(n.path)
	render, ok := RenderFns[ext]
	if !ok {
		return template.HTML(""), NotRenderableFile
	}
	return render(n, meta)
}

// RenderDir render the index html file for the directory
func (n *Node) RenderDir(m *Metadata) (template.HTML, error) {
	sb := &strings.Builder{}
	if err := dirTemplate.Execute(sb, n); err != nil {
		return template.HTML(""), err
	}
	return template.HTML(sb.String()), nil
}

func RenderMarkdown(n *Node, m *Metadata) (template.HTML, error) {
	cmd := exec.Command("kramdown", n.path)
	bytes, err := cmd.Output()
	if err != nil {
		return template.HTML(""), err
	}
	return template.HTML(bytes), nil
}

func RenderOrg(n *Node, m *Metadata) (template.HTML, error) {
	cmd := exec.Command("pandoc", "--fail-if-warnings", "-f", "org", "-t", "html", "--mathjax", "-i", filepath.Base(n.path))
	cmd.Dir = filepath.Dir(n.path)
	bytes, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("run command on %s failed: %s", n.path, string(err.(*exec.ExitError).Stderr))
		return template.HTML(""), nil
	}

	return template.HTML(bytes), nil
}

func RenderHTML(n *Node, m *Metadata) (template.HTML, error) {
	bytes, err := ioutil.ReadFile(n.path)
	if err != nil {
		return template.HTML(""), nil
	}

	return template.HTML(bytes), nil
}

package swgen

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Swgen struct {
	Source   string
	Target   string
	URLRoot  string
	Ignore   Ignore
	Template *template.Template
}

type Doc struct {
	*Swgen
	Toc  template.HTML
	Page template.HTML
	Node *Node
}

func (sw *Swgen) PageURL(n *Node) (string, error) {
	rel, err := filepath.Rel(sw.Source, n.Path)
	if err != nil {
		return "", err
	}

	url := filepath.Join(sw.URLRoot, rel)
	if !n.Info.IsDir() {
		ext := filepath.Ext(n.Info.Name())
		if _, ok := RenderFns[ext]; ok {
			url += ".html"
		}
	}
	return url, nil
}

func (sw *Swgen) Run() error {
	if err := os.MkdirAll(sw.Target, os.ModePerm); err != nil {
		return err
	}

	tree, err := sw.Scan(sw.Source)
	if err != nil {
		return err
	}

	metadata := &Metadata{}
	content := template.HTML("")
	return sw.renderAll(tree, metadata, content)
}

func (sw *Swgen) renderAll(n *Node, m *Metadata, c template.HTML) error {
	if !n.Info.IsDir() {
		return sw.render(n, m, c)
	}

	dest, err := sw.GetTargetPath(n.Path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	for _, child := range n.Children {
		err := sw.renderAll(child, m, c)
		if err != nil {
			return err
		}
	}

	return sw.renderDir(n, m, c)
}

func (sw *Swgen) copy(src string) error {
	dest, err := sw.GetTargetPath(src)
	if err != nil {
		return err
	}

	fw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fw.Close()

	fd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fw, fd)
	return err
}

var dirTemplate = template.Must(template.New("dir").Parse(`
<div>
{{.Path | .Swgen.GetRelPath}}
<ul>
{{range .Children}}
  <li>
    <a href="{{. | .Swgen.PageURL}}">{{.Info.Name}}</a>
  </li>
{{end}}
</ul>
</div>
`))

func (sw *Swgen) renderDir(n *Node, m *Metadata, c template.HTML) error {
	dest, err := sw.GetTargetPath(n.Path)
	if err != nil {
		return err
	}

	sb := &strings.Builder{}
	if err := dirTemplate.Execute(sb, n); err != nil {
		return err
	}

	fd, err := os.Create(fmt.Sprintf("%s/index.html", dest))
	if err != nil {
		return err
	}
	defer fd.Close()
	
	doc := &Doc{
		Swgen: sw,
		Toc:   c,
		Page:  template.HTML(sb.String()),
		Node:  n,
	}

	return sw.Template.Execute(fd, doc)
}

func(sw *Swgen) GetRelPath(path string) (string, error) {
	return filepath.Rel(sw.Source, path)
}

func (sw *Swgen) GetTargetPath(path string) (string, error) {
	rel, err := sw.GetRelPath(path)
	if err != nil {
		return "", err
	}

	return filepath.Join(sw.Target, rel), nil
}

func (sw *Swgen) render(n *Node, m *Metadata, c template.HTML) error {
	dest, err := sw.GetTargetPath(n.Path)
	if err != nil {
		return err
	}

	html, err := n.Render(m)
	if err != nil {
		return err
	}

	log.Printf("render %s to %s.html", n.Path, dest)
	fd, err := os.Create(fmt.Sprintf("%s.html", dest))
	if err != nil {
		return err
	}
	defer fd.Close()

	doc := &Doc{
		Swgen: sw,
		Toc:   c,
		Page:  html,
		Node:  n,
	}

	return sw.Template.Execute(fd, doc)
}

func (sw *Swgen) Scan(root string) (*Node, error) {
	// the root must be a directory
	info, err := os.Lstat(root)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("root %s must be a directory", root)
	}

	home := &Node{
		Swgen:    sw,
		Info:     info,
		Path:     root,
		Children: []*Node{},
	}

	return sw.scan(root, info, home)
}

func (sw *Swgen) scan(path string, info os.FileInfo, home *Node) (*Node, error) {
	n := &Node{
		Swgen:    sw,
		Info:     info,
		Path:     path,
		Children: []*Node{},
		Home:     home,
	}

	if info.IsDir() {
		dir, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		children, err := dir.Readdir(0)
		if err != nil {
			return nil, err
		}

		for _, child := range children {
			path := filepath.Join(path, child.Name())
			if sw.Ignore.Ignore(path) {
				log.Printf("ignore path %s", path)
				continue
			}

			if !child.IsDir() {
				ext := filepath.Ext(child.Name())
				if _, ok := RenderFns[ext]; !ok {
					sw.copy(path)
					continue
				}
			}

			subNode, err := sw.scan(path, child, home)
			if err != nil {
				return nil, err
			}

			subNode.Up = n
			n.Children = append(n.Children, subNode)
		}

		// update prev/next link
		for i := 0; i < len(n.Children); i++ {
			if i > 0 {
				n.Children[i-1].Next = n.Children[i]
			}

			if i < len(n.Children)-1 {
				n.Children[i+1].Prev = n.Children[i]
			}
		}
	}

	return n, nil
}

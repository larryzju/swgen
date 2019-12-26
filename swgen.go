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

// Swgen is the main structure to scan source directory and render pages to target directory
type Swgen struct {
	Source   string
	Target   string
	URLRoot  string
	Ignore   Ignore
	Force    bool
	Template *template.Template
}

// Doc is the virtual page object to render
type Doc struct {
	Toc  template.HTML
	Page template.HTML
	*Node
}

// Run scans source directory and render pages to output directory
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
	dest := sw.MustGetTargetPath(n.path)

	if n.Info.IsDir() {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return err
		}

		dest := fmt.Sprintf("%s/index.html", sw.MustGetTargetPath(n.path))
		html, err := n.RenderDir(m)
		if err != nil {
			return err
		}
		if err := sw.render(dest, n, c, html); err != nil {
			return err
		}

		for _, child := range n.Children {
			err := sw.renderAll(child, m, c)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// regular files
	suffix := filepath.Ext(dest)
	if !(strings.EqualFold(suffix, "html") || strings.EqualFold(suffix, "htm")) {
		dest = fmt.Sprintf("%s.html", dest)
	}

	// if the target exists and force flag is not enable, then skip the generate
	destInfo, err := os.Stat(dest)
	if err == nil && destInfo.ModTime().After(n.Info.ModTime()) && !sw.Force {
		log.Printf("skip existed file %s", dest)
		return nil
	}

	html, err := n.Render(m)
	if err != nil {
		return err
	}

	return sw.render(dest, n, c, html)
}

func (sw *Swgen) render(dest string, n *Node, c, html template.HTML) error {
	log.Printf("render %s to %s", n.path, dest)
	fd, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fd.Close()

	doc := &Doc{
		Toc:  c,
		Page: html,
		Node: n,
	}

	return sw.Template.Execute(fd, doc)
}

func (sw *Swgen) copy(src string) error {
	dest := sw.MustGetTargetPath(src)
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
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

// MustGetRelPath get the relative path
func (sw *Swgen) MustGetRelPath(path string) string {
	rel, err := filepath.Rel(sw.Source, path)
	if err != nil {
		log.Panicf("rel %s %s failed", sw.Source, path)
	}
	return rel
}

// GetTargetPath get the target path
func (sw *Swgen) MustGetTargetPath(path string) string {
	rel := sw.MustGetRelPath(path)
	return filepath.Join(sw.Target, rel)
}

// Scan the source directory and return nodes tree
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
		path:     root,
		Children: []*Node{},
	}

	return sw.scan(root, info, home)
}

func (sw *Swgen) scan(path string, info os.FileInfo, home *Node) (*Node, error) {
	n := &Node{
		Swgen:    sw,
		Info:     info,
		path:     path,
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
			if sw.Ignore.Ignore(sw.MustGetRelPath(path)) {
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

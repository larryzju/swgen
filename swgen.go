package swgen

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/larryzju/swgen/node"
)

const (
	contentHtmlTempalteLiteral = `
	<span>{{.Title}}</span>
	<ul>{{range .Children}}<li>{{.Link}}</li>{{end}}</ul>
	`
)

var (
	DiaryTemplate *template.Template
)

func init() {
	DiaryTemplate = template.Must(template.ParseGlob("template/*.html"))
}

type HTMLPage struct {
	Title   string
	Content template.HTML
	Body    template.HTML
}

type Dir interface {
	Generate(root string, content template.HTML) error
	Children() []node.HTMLNode
}

type DirNode struct {
	rel      string
	path     string
	info     os.FileInfo
	children []node.Node
}

func NewDirNode(root string, srcdir string, ignore Ignore) (e node.Node, err error) {
	info, err := os.Stat(srcdir)
	if err != nil {
		return
	}

	rel, err := filepath.Rel(root, srcdir)
	if err != nil {
		return
	}

	dir, err := os.Open(srcdir)
	if err != nil {
		return
	}
	defer dir.Close()

	infos, err := dir.Readdir(0)
	if err != nil {
		return
	}

	d := &DirNode{
		rel:  rel,
		path: srcdir,
		info: info,
	}

	for _, sub := range infos {
		var c node.Node
		pagePath := path.Join(srcdir, sub.Name())
		if sub.IsDir() {
			c, err = NewDirNode(root, pagePath, ignore)
		} else {
			fn := getNewNodeFn(sub)
			c, err = fn(root, pagePath)
		}

		if err != nil {
			return nil, err
		}

		if ignore.Ignore(c.Rel()) {
			log.Println("ignore", pagePath)
			continue
		}

		d.children = append(d.children, c)
	}

	return d, nil
}

func getNewNodeFn(info os.FileInfo) node.NewFn {
	switch path.Ext(info.Name()) {
	case ".org":
		return node.NewOrgNode
	// case "md":
	// 	page.NewMarkdownElemnt
	default:
		return node.NewRawNode
	}
}

func (d *DirNode) Rel() string {
	return d.rel
}

func (d *DirNode) Title() string {
	return d.info.Name()
}

func (d *DirNode) LastUpdate() time.Time {
	return d.info.ModTime()
}

func (d *DirNode) String() string {
	return d.Rel()
}

func (d *DirNode) Children() []node.HTMLNode {
	nodes := []node.HTMLNode{}
	for _, c := range d.children {
		if htmlNode, ok := c.(node.HTMLNode); ok {
			nodes = append(nodes, htmlNode)
		}
	}
	return nodes
}

func (d *DirNode) Link(root string) string {
	return path.Join(root, d.Rel())
}

func (d *DirNode) Content() (template.HTML, error) {
	return template.HTML("TODO"), nil
}

func (d *DirNode) Navbar(root string) (template.HTML, error) {
	t := `
	  <a href='{{.selfLink}}'>{{.title}}</a>
	  <ul>
	  {{range .subLinks}}
	    <li>{{.}}</li>
	  {{end}}
	  </ul>
	`

	subLinks := []template.HTML{}
	for _, c := range d.children {
		if dir, ok := c.(*DirNode); ok {
			if len(dir.children) == 0 {
				continue
			}
			subNavbar, err := dir.Navbar(root)
			if err != nil {
				return "", err
			}
			subLinks = append(subLinks, subNavbar)
		} else if html, ok := c.(node.HTMLNode); ok {
			subLinks = append(subLinks, template.HTML(fmt.Sprintf("<a href='%s'>%s</a>", html.Link(root), c.Title())))
		}
	}

	v := map[string]interface{}{
		"selfLink": d.Link(root),
		"title":    d.Title(),
		"subLinks": subLinks,
	}

	b := &strings.Builder{}
	tmpl := template.Must(template.New("navbar").Parse(t))
	err := tmpl.Execute(b, v)
	if err != nil {
		return "", err
	}
	return template.HTML(b.String()), nil
}

func (d *DirNode) Generate(root string, navbar template.HTML) error {
	// create directory if not exists
	dst := path.Join(root, d.Rel())
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	for _, c := range d.children {
		switch c.(type) {
		case Dir:
			c.(*DirNode).Generate(root, navbar)

		case node.RawNode:
			raw := c.(node.RawNode)
			log.Printf("open %s to write raw\n", path.Join(root, c.Rel()))
			w, err := os.Create(path.Join(root, c.Rel()))
			if err != nil {
				return err
			}
			defer w.Close()

			r, err := raw.Reader()
			if err != nil {
				return err
			}
			defer r.Close()

			_, err = io.Copy(w, r)
			if err != nil {
				return err
			}

		case node.HTMLNode:
			html := c.(node.HTMLNode)
			log.Printf("open %s to write html\n", path.Join(root, c.Rel()))
			w, err := os.Create(path.Join(root, c.Rel()))
			if err != nil {
				return err
			}
			defer w.Close()

			body, err := html.Content()
			if err != nil {
				return err
			}

			page := HTMLPage{
				Title:   c.Title(),
				Content: navbar,
				Body:    body,
			}

			err = DiaryTemplate.Execute(w, page)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

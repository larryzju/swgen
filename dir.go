package swgen

import (
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/larryzju/swgen/node"
)

type Dir struct {
	rel   string
	path  string
	info  os.FileInfo
	nodes []node.Node
}

type Source struct {
	root   string
	path   string
	ignore Ignore
}

func (s *Source) Root() string { return s.root }
func (s *Source) Path() string { return s.path }
func (s *Source) Ext() string  { return path.Ext(s.path) }
func (s *Source) Rel() string {
	rel, err := filepath.Rel(s.root, s.path)
	if err != nil {
		panic(err)
	}
	return rel
}

type Target struct {
	root    string
	urlRoot string
}

func (t *Target) Root() string    { return t.root }
func (t *Target) URLRoot() string { return t.urlRoot }

type Metadata struct {
	navigator  string
	buildTime  time.Time
	gitVersion string
}

func (m *Metadata) Navigator() template.HTML { return template.HTML(m.navigator) }
func (m *Metadata) BuildTime() time.Time     { return m.buildTime }
func (m *Metadata) GitVersion() string       { return m.gitVersion }

func newDir(s node.Source) (node node.Node, err error) {
	return NewDir(s.(*Source))
}

func NewDir(s *Source) (d *Dir, err error) {
	info, err := os.Stat(s.path)
	if err != nil {
		return
	}

	rel, err := filepath.Rel(s.root, s.path)
	if err != nil {
		return
	}

	dir, err := os.Open(s.path)
	if err != nil {
		return
	}
	defer dir.Close()

	infos, err := dir.Readdir(0)
	if err != nil {
		return
	}

	d = &Dir{
		rel:  rel,
		path: s.path,
		info: info,
	}

	for _, c := range infos {
		p := path.Join(s.path, c.Name())
		s := Source{root: s.root, path: p, ignore: s.ignore}
		if s.ignore.Ignore(s.Rel()) {
			log.Println("ignore", p)
			continue
		}

		newFn := node.New
		if c.IsDir() {
			newFn = newDir
		}

		node, err := newFn(&s)
		if err != nil {
			return nil, err
		}
		d.nodes = append(d.nodes, node)
	}

	return d, nil
}

func (d *Dir) Rel() string {
	return d.rel
}

func (d *Dir) Title() string {
	return d.info.Name()
}

func (d *Dir) LastUpdate() time.Time {
	return d.info.ModTime()
}

func (d *Dir) String() string {
	return d.Rel()
}

func (d *Dir) Flush(m node.Metadata, t node.Target) (err error) {
	dst := path.Join(t.Root(), d.Rel())
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	for _, c := range d.nodes {
		err = c.Flush(m, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Dir) generateNavigator(t node.Target) (navigator string, err error) {
	links := []template.HTML{}
	for _, n := range d.nodes {
		switch n.(type) {
		case *node.Raw:
			continue
		case *Dir:
			link, err := n.(*Dir).generateNavigator(t)
			if err != nil {
				return "", err
			}

			if link == "" {
				continue
			}
			links = append(links, template.HTML("<span>"+n.Title()+"</span>"+link))
		default:
			nt := &struct {
				node.Node
				node.Target
			}{n, t}
			tmpl := template.Must(template.New("item").Parse(`<a href='{{.URLRoot}}/{{.Rel}}'>{{.Title}}</a>`))
			sb := &strings.Builder{}
			tmpl.Execute(sb, nt)
			links = append(links, template.HTML(sb.String()))
		}
	}

	if len(links) == 0 {
		return "", nil
	}

	sb := &strings.Builder{}
	listTemplate := template.Must(template.New("list").Parse(`<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>`))
	listTemplate.Execute(sb, links)
	return sb.String(), nil
}
package swgen

import (
	"log"
	"os"
	"path"
	"path/filepath"
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

func (m *Metadata) NavigatorURL() string { return m.navigator }
func (m *Metadata) BuildTime() time.Time { return m.buildTime }
func (m *Metadata) GitVersion() string   { return m.gitVersion }

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
		dstPath := path.Join(t.Root(), c.Rel())
		dstInfo, err := os.Stat(dstPath)
		if err == nil && !dstInfo.IsDir() && dstInfo.ModTime().After(d.LastUpdate()) {
			continue
		}
		err = c.Flush(m, t)
		if err != nil {
			return err
		}
	}

	return nil
}

type Navigator struct {
	Title    string       `json:"title"`
	Link     string       `json:"link,omitempty"`
	Children []*Navigator `json:"children,omitempty"`
}

func (d *Dir) generateNavigator(t node.Target) (nav *Navigator, err error) {
	nav = &Navigator{Title: d.Title()}
	for _, n := range d.nodes {
		switch n.(type) {
		case *node.Raw:
			continue
		case *Dir:
			nav0, err := n.(*Dir).generateNavigator(t)
			if err != nil {
				return nil, err
			}

			if len(nav0.Children) == 0 {
				continue
			}
			nav.Children = append(nav.Children, nav0)
		default:
			nav0 := &Navigator{Title: n.Title(), Link: path.Join(t.URLRoot(), n.Rel())}
			nav.Children = append(nav.Children, nav0)
		}
	}
	return
}

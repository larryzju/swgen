package swgen

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/larryzju/swgen/page"
)

type Page interface {
	Rel() string
	Title() string
	LastUpdate() time.Time
	Write(w io.Writer) (int64, error)
	ContentHTML() string
	HTML() string
}

type DirPage struct {
	rel      string
	path     string
	info     os.FileInfo
	children []Page
}

func NewDirPage(root string, srcdir string) (d *DirPage, err error) {
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

	d = &DirPage{
		rel:  rel,
		path: srcdir,
		info: info,
	}

	for _, info := range infos {
		pagePath := path.Join(srcdir, info.Name())
		p, err := generatePage(root, pagePath, info)
		if err != nil {
			return nil, err
		}
		d.children = append(d.children, p)
	}

	return
}

func generatePage(root, path string, info os.FileInfo) (p Page, err error) {
	if info.IsDir() {
		p, err = NewDirPage(root, path)
	} else {
		p, err = page.NewRawPage(root, path)
	}
	return
}

func (d *DirPage) Rel() string {
	return d.rel
}

func (d *DirPage) Children() []Page {
	return d.children
}

func (d *DirPage) Title() string {
	return d.info.Name()
}

func (d *DirPage) LastUpdate() time.Time {
	return d.info.ModTime()
}

func (d *DirPage) Write(w io.Writer) (written int64, err error) {
	b := strings.NewReader("TODO")
	return io.Copy(w, b)
}

func (d *DirPage) String() string {
	return d.Rel()
}

func (d *DirPage) ContentHTML() string {
	w := &strings.Builder{}
	err := contentHtmlTemplate.Execute(w, d)
	if err != nil {
		log.Panic(err)
	}
	return w.String()
}

func (d *DirPage) HTML() string {
	return "CONTENT"
}

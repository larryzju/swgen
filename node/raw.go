package node

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type Raw struct {
	rel  string
	path string
	info os.FileInfo
}

func NewRawNode(root string, src string) (page Node, err error) {
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

	return &Raw{
		rel:  rel,
		path: path,
		info: info,
	}, nil
}

func (p *Raw) Title() string {
	return p.info.Name()
}

func (p *Raw) Rel() string {
	return p.rel
}

func (p *Raw) LastUpdate() time.Time {
	return p.info.ModTime()
}

func (p *Raw) String() string {
	return p.Rel()
}

func (p *Raw) Reader() (io.ReadCloser, error) {
	return os.Open(p.path)
}

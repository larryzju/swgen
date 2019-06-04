package node

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Raw struct {
	rel  string
	path string
	info os.FileInfo
}

func NewRaw(s Source) (page Node, err error) {
	root := s.Root()
	path := s.Path()
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	rel, err := filepath.Rel(root, path)
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

func (p *Raw) Flush(m Metadata, t Target) (err error) {
	dstpath := path.Join(t.Root(), p.Rel())
	dst, err := os.Create(dstpath)
	if err != nil {
		return err
	}
	defer dst.Close()

	src, err := os.Open(p.path)
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	return
}

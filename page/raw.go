package page

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type RawPage struct {
	rel  string
	path string
	file *os.File
	info os.FileInfo
}

func NewRawPage(root string, src string) (page *RawPage, err error) {
	path, err := filepath.Abs(src)
	if err != nil {
		return
	}

	file, err := os.Open(src)
	if err != nil {
		return
	}

	info, err := file.Stat()
	if err != nil {
		return
	}

	rel, err := filepath.Rel(root, src)
	if err != nil {
		return
	}

	return &RawPage{
		rel:  rel,
		path: path,
		file: file,
		info: info,
	}, nil
}

func (p *RawPage) Title() string {
	return p.info.Name()
}

func (p *RawPage) Rel() string {
	return p.rel
}

func (p *RawPage) LastUpdate() time.Time {
	return p.info.ModTime()
}

func (p *RawPage) Write(w io.Writer) (written int64, err error) {
	return io.Copy(w, p.file)
}

func (p *RawPage) String() string {
	return p.Rel()
}

func (p *RawPage) HTML() string {
	return "BODY"
}

func (p *RawPage) ContentHTML() string {
	return "<li>" + p.Title() + "</li>"
}

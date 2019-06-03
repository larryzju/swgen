package swgen

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/gobwas/glob"
)

type Ignore interface {
	Ignore(path string) bool
}

type BasicIgnore struct {
	rules []glob.Glob
}

func (igf *BasicIgnore) Ignore(path string) bool {
	for _, r := range igf.rules {
		if r.Match(path) {
			return true
		}
	}

	return false
}

func NewBasicIgnore(r io.Reader) *BasicIgnore {
	scan := bufio.NewScanner(r)
	bi := &BasicIgnore{}
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		bi.rules = append(bi.rules, glob.MustCompile(line, os.PathSeparator))
	}
	return bi
}

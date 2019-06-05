package node

import (
	"html/template"
	"time"
)

type Node interface {
	Rel() string
	Title() string
	LastUpdate() time.Time
	Flush(Metadata, Target) error
}

type Source interface {
	Root() string
	Path() string
	Ext() string
}

type Target interface {
	URLRoot() string
	Root() string
}

type Metadata interface {
	Navigator() template.HTML
	BuildTime() time.Time
	GitVersion() string
}

type NewFn func(Source) (Node, error)

func New(s Source) (Node, error) {
	switch s.Ext() {
	case ".org":
		return NewOrg(s)
	case ".md":
		return NewMarkdown(s)
	default:
		return NewRaw(s)
	}
}

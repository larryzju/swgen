package node

import (
	"html/template"
	"io"
	"time"
)

type Node interface {
	Rel() string
	Title() string
	LastUpdate() time.Time
}

type HTMLNode interface {
	Node
	Link(root string) string
	Content() (template.HTML, error)
}

type RawNode interface {
	Node
	Reader() (io.ReadCloser, error)
}

type NewFn func(root string, path string) (Node, error)

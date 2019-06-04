package node

import (
	"html/template"
)

var (
	htmlTemplate *template.Template
)

func init() {
	htmlTemplate = template.Must(template.ParseFiles("template/page.html"))
}

type HTMLData interface {
	Node
	Metadata
	Body() template.HTML
}

type HTMLDataWrapper struct {
	Node
	Metadata
	Body template.HTML
}

func wrapHTMLData(n Node, m Metadata, b template.HTML) *HTMLDataWrapper {
	return &HTMLDataWrapper{Node: n, Metadata: m, Body: b}
}

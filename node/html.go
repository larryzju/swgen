package node

import (
	"html/template"
)

const (
	PageTemplateLiteral = `
	<!DOCTYPE html>
	<html>
		<head>
			<title>{{.Title}}</title>
			<link rel="stylesheet" href="{{.URLRoot}}/resources/css/main.css" />
			<script src="{{.URLRoot}}/resources/js/jquery-3.4.1.min.js"></script>
			<script src="{{.URLRoot}}/resources/js/swgen.js"></script>
		</head>
		<body>
			<nav id='content'>
				{{.Navigator}}
			</nav>
			<section class="main-article-area">
				<div id='main'>
					{{.Body}}
				</div>
			</section>
		</body>
	</html>
	`
)

var (
	htmlTemplate *template.Template
)

func init() {
	htmlTemplate = template.Must(template.New("page").Parse(PageTemplateLiteral))
}

type HTMLData interface {
	Title() string
	URLRoot() string
	Navigator() template.HTML
	Body() template.HTML
}

type HTMLDataWrapper struct {
	Node
	Metadata
	Target
	Body template.HTML
}

func wrapHTMLData(n Node, m Metadata, t Target, b template.HTML) *HTMLDataWrapper {
	return &HTMLDataWrapper{Node: n, Target: t, Metadata: m, Body: b}
}

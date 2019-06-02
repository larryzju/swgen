package swgen

import (
	template "text/template"
)

const (
	contentHtmlTempalteLiteral = `
	<span>
		<a href='{{.Rel}}'>{{.Title}}</a>
		<ul>
			{{range .Children}}
			<li><a href='{{.Rel}}'>{{.ContentHTML}}</a></li>{{end}}
		</ul>
	</span>
	`
)

var (
	contentHtmlTemplate *template.Template
)

func init() {
	contentHtmlTemplate = template.Must(
		template.New("content").Parse(contentHtmlTempalteLiteral))
}

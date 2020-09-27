package tmpl

const FileHTML = `
{{ define "content" }}
<img src="{{ .FileServerURL }}/{{ .Filename }}" class="file">
{{ end }}
`

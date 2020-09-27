package templates

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/packago/config"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tullo/imgasm.com/templates/tmpl"
)

var tmpls map[string]*template.Template = make(map[string]*template.Template)

func init() {
	if config.File().GetString("environment") == "development" {
		minifyCSS()
	}
	layout := minifyHTML(tmpl.LayoutHTML)

	tmpls["index.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["index.html"].New("index").Parse(minifyHTML(tmpl.IndexHTML)))

	tmpls["file.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["file.html"].New("file").Parse(minifyHTML(tmpl.FileHTML)))

	tmpls["about.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["about.html"].New("about").Parse(minifyHTML(tmpl.AboutHTML)))
	tmpls["privacy-policy.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["privacy-policy.html"].New("privacy-policy").Parse(minifyHTML(tmpl.PrivacyPolicyHTML)))

	tmpls["rate-limit.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["rate-limit.html"].New("rate-limit").Parse(minifyHTML(tmpl.RateLimitHTML)))
	tmpls["not-found.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["not-found.html"].New("not-found").Parse(minifyHTML(tmpl.NotFoundHTML)))
}

func Render(wr io.Writer, template string, data interface{}) {
	err := tmpls[template].Execute(wr, data)
	if err != nil {
		log.Panicf("[ERROR] Error rendering %s: %s\n", template, err)
	}
}

func minifyHTML(input string) string {
	min := minify.New()
	min.Add("text/html", &html.Minifier{})
	str, err := min.String("text/html", input)
	if err != nil {
		panic(err)
	}
	return str
}

func minifyCSS() {
	filepaths := []string{"./ui/static/css/main.css"}
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	var sb strings.Builder
	for _, filepath := range filepaths {
		css, err := ioutil.ReadFile(filepath)
		if err != nil {
			panic(err)
		}
		min, err := m.String("text/css", string(css))
		if err != nil {
			panic(err)
		}
		sb.WriteString(min)
	}
	err := ioutil.WriteFile("./ui/static/main.min.css", []byte(sb.String()), 0644)
	if err != nil {
		panic(err)
	}
}

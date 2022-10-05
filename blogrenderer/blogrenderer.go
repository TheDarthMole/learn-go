package blogrenderer

import (
	"embed"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"html/template"
	"io"
	"strings"
)

//go:embed "templates/*"
var templates embed.FS

type PostRenderer struct {
	tmpl *template.Template
}

type Post struct {
	Title, Description, Body string
	Tags                     []string
}

func NewPostRenderer() (*PostRenderer, error) {
	tmpl, err := template.New("postRenderer").Funcs(
		template.FuncMap{}).ParseFS(templates, "templates/*.gohtml")
	return &PostRenderer{tmpl: tmpl}, err
}

func (p *Post) SanitisedTitle() string {
	return strings.ToLower(strings.Replace(p.Title, " ", "-", -1))
}

func (r *PostRenderer) Render(w io.Writer, p Post) error {
	return r.tmpl.ExecuteTemplate(w, "post.gohtml", newPostVM(p))
}

func (r *PostRenderer) RenderIndex(w io.Writer, p []Post) error {
	return r.tmpl.ExecuteTemplate(w, "index.gohtml", p)
}

type postViewModel struct {
	Post     *Post
	HTMLBody template.HTML
}

func newParser() *parser.Parser {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	return parser.NewWithExtensions(extensions)
}

func newPostVM(p Post) postViewModel {
	vm := postViewModel{Post: &p}
	md := markdown.ToHTML([]byte(p.Body), newParser(), nil)
	vm.HTMLBody = template.HTML(md)
	return vm
}

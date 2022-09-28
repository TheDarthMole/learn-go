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
	tmpl     *template.Template
	mdParser *parser.Parser
}

type Post struct {
	Title, Description, Body string
	Tags                     []string
}

func NewPostRenderer() (*PostRenderer, error) {
	tmpl, err := template.New("postRenderer").Funcs(
		template.FuncMap{}).ParseFS(templates, "templates/*.gohtml")

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	return &PostRenderer{tmpl: tmpl, mdParser: parser}, err
}

func (p *Post) SanitisedTitle() string {
	return strings.ToLower(strings.Replace(p.Title, " ", "-", -1))
}

func (r *PostRenderer) Render(w io.Writer, p Post) error {
	return r.tmpl.ExecuteTemplate(w, "post.gohtml", newPostVM(p, r))
}

func (r *PostRenderer) RenderIndex(w io.Writer, p []Post) error {
	return r.tmpl.ExecuteTemplate(w, "index.gohtml", p)
}

type postViewModel struct {
	Post     *Post
	HTMLBody template.HTML
}

func newPostVM(p Post, r *PostRenderer) postViewModel {
	vm := postViewModel{Post: &p}
	vm.HTMLBody = template.HTML(markdown.ToHTML([]byte(p.Body), r.mdParser, nil))
	return vm
}

package blogrenderer_test

import (
	"bytes"
	approvals "github.com/approvals/go-approval-tests"
	"io"
	"learn-go/blogrenderer"
	"testing"
)

func TestRender(t *testing.T) {

	var aPost = blogrenderer.Post{
		Title:       "hello world",
		Description: "This is a description",
		Body:        "# Title\n\nbody here",
		Tags:        []string{"go", "tdd"},
	}

	postRenderer, err := blogrenderer.NewPostRenderer()

	if err != nil {
		t.Fatal(err)
	}

	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}

		if err = postRenderer.Render(&buf, aPost); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})

	t.Run("an index page with our posts is rendered", func(t *testing.T) {
		buf := bytes.Buffer{}
		posts := []blogrenderer.Post{
			{Title: "Hello world 1"},
			{Title: "Hello world 2"},
		}

		if err = postRenderer.RenderIndex(&buf, posts); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}

func BenchmarkReader(b *testing.B) {
	var aPost = blogrenderer.Post{
		Title:       "Hello World!",
		Body:        "This is a post",
		Description: "This is a description",
		Tags:        []string{"go", "tdd"},
	}

	b.ResetTimer()

	tmpl, err := blogrenderer.NewPostRenderer()
	if err != nil {
		b.Fatalf("error getting a valid template from templating engine, err: %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		err = tmpl.Render(io.Discard, aPost)
		if err != nil {
			b.Fatal(err)
		}
	}
}

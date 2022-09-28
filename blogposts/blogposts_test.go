package blogposts_test

import (
	"errors"
	"io/fs"
	"learn-go/blogposts"
	"reflect"
	"testing"
	"testing/fstest"
)

type StubFailingFS struct {
}

const (
	firstPost = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
hello,
world`
	secondPost = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker
---
B
L
M`
)

func (fs StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no! I always fail")
}

func TestNewBlogPost(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstPost)},
		"hello-world2.md": {Data: []byte(secondPost)},
	}

	posts, err := blogposts.NewPostsFromFS(fs)

	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != len(fs) {
		t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
	}

	got := posts[0]
	want := blogposts.Post{
		Title:       "Post 1",
		Description: "Description 1",
		Tags:        []string{"tdd", "go"},
		Body:        "hello,\nworld",
	}

	assertPost(t, got, want)
}

func assertPost(t *testing.T, got, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

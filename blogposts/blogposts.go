package blogposts

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"strings"
	"testing/fstest"
)

const (
	titleTag = "Title: "
	descTag  = "Description: "
	tagsTag  = "Tags: "
)

type Post struct {
	Title, Description, Body string
	Tags                     []string
}

func NewPostsFromFS(filesystem fstest.MapFS) ([]Post, error) {

	dir, err := fs.ReadDir(filesystem, ".")

	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, file := range dir {
		post, err := getPost(filesystem, file.Name())
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(filesystem fstest.MapFS, fileName string) (Post, error) {
	postFile, err := filesystem.Open(fileName)

	if err != nil {
		return Post{}, err
	}
	defer func(postFile fs.File) {
		err = postFile.Close()
		if err != nil {
			log.Fatalf("error closing file: %+v", err)
		}
	}(postFile)

	return newPost(postFile)
}

func newPost(file io.Reader) (Post, error) {
	scanner := bufio.NewScanner(file)

	readMetaLine := func(tag string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tag)
	}

	return Post{
		Title:       readMetaLine(titleTag),
		Description: readMetaLine(descTag),
		Tags:        strings.Split(readMetaLine(tagsTag), ", "),
		Body:        readBody(*scanner),
	}, nil
}

func readBody(scanner bufio.Scanner) string {
	buf := bytes.Buffer{}
	scanner.Scan() // discard the `---` as we don't need it
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

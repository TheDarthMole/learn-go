package poker

import (
	"io"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	testTape := &tape{file}

	if _, err := testTape.Write([]byte("abc")); err != nil {
		t.Errorf("not expecting error when writing, got: %q", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Errorf("not expecting seek error, got: %q", err)
	}
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

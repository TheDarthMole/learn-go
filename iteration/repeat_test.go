package iteration

import (
	"testing"
)

func AssertEqualsString(t *testing.T, expected, got string) {
	if got != expected {
		t.Errorf("Expected %q got %q", expected, got)
	}
}

func TestRepeat(t *testing.T) {

	t.Run("Test repeats 6 times", func(t *testing.T) {
		repeated := Repeat("A", 6)
		expected := "AAAAAA"
		AssertEqualsString(t, expected, repeated)
	})

	t.Run("Test repeats 0 times", func(t *testing.T) {
		repeated := Repeat("A", 0)
		expected := ""
		AssertEqualsString(t, expected, repeated)
	})

}

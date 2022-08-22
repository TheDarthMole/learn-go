package concurrency

import (
	"reflect"
	"testing"
	"time"
)

func mockWebsiteChecker(url string) bool {
	if url == "waat://furhurterwe.geds" {
		return false
	}
	return true
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"https://google.com",
		"https://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"https://google.com":          true,
		"https://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":     false,
	}

	t.Run("test no delay", func(t *testing.T) {
		got := CheckWebsites(mockWebsiteChecker, websites)

		if !reflect.DeepEqual(want, got) {
			t.Fatalf("wanted %v, got %v", want, got)
		}
	})

	t.Run("test 20ms delay", func(t *testing.T) {
		got := CheckWebsites(slowStubWebsiteChecker, websites)

		if !reflect.DeepEqual(want, got) {
			t.Fatalf("wanted %v, got %v", want, got)
		}
	})
}

func slowStubWebsiteChecker(url string) bool {
	time.Sleep(20 * time.Millisecond)
	return mockWebsiteChecker(url)
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)

	for i := 0; i < len(urls); i++ {
		urls[i] = "a url"
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CheckWebsites(slowStubWebsiteChecker, urls)
	}
}

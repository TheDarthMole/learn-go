package racer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}

func TestRacer(t *testing.T) {
	slowServer := makeDelayedServer(20 * time.Millisecond)
	fastServer := makeDelayedServer(0)

	// `defer` calls the function at the end of the containing function
	defer slowServer.Close()
	defer fastServer.Close()

	t.Run("test fast vs slow server", func(t *testing.T) {
		want := fastServer.URL
		got := Racer(slowServer.URL, fastServer.URL)

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

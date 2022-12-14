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

	testTimeout := 1 * time.Second

	// `defer` calls the function at the end of the containing function
	defer slowServer.Close()
	defer fastServer.Close()

	t.Run("test fast vs slow server", func(t *testing.T) {
		want := fastServer.URL
		got, err := Racer(slowServer.URL, fastServer.URL)

		if err != nil {
			t.Errorf("got an error when not expecting one")
		}

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("test that testTimeout causes error", func(t *testing.T) {
		_, err := ConfigurableRacer("http://localhost:1", "http://localhost:2", testTimeout)

		if err == nil {
			t.Errorf("expected to get error %s but got nil ", err)
		}
	})
}

package racer

import (
	"fmt"
	"net/http"
	"time"
)

func websitePing(url string) chan struct{} {
	ch := make(chan struct{})

	go func() {
		_, err := http.Get(url)
		if err == nil {
			close(ch)
		}
	}()
	return ch
}

var timeout = 10 * time.Second

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
	// We don't care for the return value of websitePing
	// We just care that it closes the channel, which triggers
	// the `select` statement and executes the code below it.
	// First to return is executed

	select {
	case <-websitePing(a):
		return a, nil
	case <-websitePing(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out after trying to get '%s' and '%s'", a, b)
	}
}

func Racer(a, b string) (winner string, error error) {
	return ConfigurableRacer(a, b, timeout)
}

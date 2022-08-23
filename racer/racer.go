package racer

import (
	"net/http"
	"time"
)

func websitePinger(url string, out chan time.Duration) {
	start := time.Now()
	http.Get(url)
	out <- time.Since(start)
}

func Racer(a, b string) string {
	chanA := make(chan time.Duration)
	chanB := make(chan time.Duration)

	go func() { websitePinger(a, chanA) }()
	go func() { websitePinger(b, chanB) }()

	if <-chanA < <-chanB {
		return a
	}
	return b
}

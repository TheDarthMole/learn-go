package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

type SpyCountdownOperator struct {
	Calls []string
}

func (s *SpyCountdownOperator) Sleep(duration time.Duration) {
	s.Calls = append(s.Calls, "sleep")
}

func (s *SpyCountdownOperator) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, "write")
	return
}

const (
	write = "write"
	sleep = "sleep"
)

func TestCountdown(t *testing.T) {

	t.Run("sleep before every print", func(t *testing.T) {
		spyCountDownOperator := &SpyCountdownOperator{}
		Countdown(spyCountDownOperator, spyCountDownOperator)

		want := []string{
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}

		if !reflect.DeepEqual(want, spyCountDownOperator.Calls) {
			t.Errorf("got %v want %v", spyCountDownOperator.Calls, want)
		}
	})

	t.Run("prints 3 to Go!", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		Countdown(buffer, &SpyCountdownOperator{})

		got := buffer.String()
		want := `3
2
1
Go!`
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

}

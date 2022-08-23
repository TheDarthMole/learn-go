package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Sleeper interface {
	Sleep(duration time.Duration)
}

type DefaultSleeper struct{}
type DefaultOut struct{}

func (d *DefaultSleeper) Sleep(ms time.Duration) {
	time.Sleep(ms * time.Second)
}

func Countdown(out io.Writer, sleeper Sleeper) {
	for i := 3; i > 0; i-- {
		_, err := fmt.Fprintln(out, i)
		if err != nil {
			return
		}
		sleeper.Sleep(1 * time.Second)
	}
	_, err := fmt.Fprint(out, "Go!")
	if err != nil {
		return
	}
}

var sleeper = &DefaultSleeper{}

func main() {
	Countdown(os.Stdout, sleeper)
	// Countdown()
}

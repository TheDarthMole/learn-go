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

func (d *DefaultOut) Write(out io.Writer) {

}

func Countdown(out io.Writer, sleeper Sleeper) {
	for i := 3; i > 0; i-- {
		fmt.Fprintln(out, i)
		sleeper.Sleep(1 * time.Second)
	}
	fmt.Fprint(out, "Go!")
}

var sleeper = &DefaultSleeper{}

func main() {
	Countdown(os.Stdout, sleeper)
	// Countdown()
}

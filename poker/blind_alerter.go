package poker

import (
	"fmt"
	"io"
	"log"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int, to io.Writer)
}

type BlindAlerterFunc func(duration time.Duration, amount int, to io.Writer)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	a(duration, amount, to)
}

func StdOutAlerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		_, err := fmt.Fprintf(to, "Blind is now %d\n", amount)
		if err != nil {
			log.Fatalf("could not log to dummyStdOut")
		}
	})
}

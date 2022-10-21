package poker

import (
	"fmt"
	"log"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type BlindAlerterFunc func(duration time.Duration, amount int)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		_, err := fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
		if err != nil {
			log.Fatalf("could not log to stdout")
		}
	})
}

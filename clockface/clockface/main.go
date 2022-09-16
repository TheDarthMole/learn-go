package main

import (
	"learn-go/clockface"
	"os"
	"time"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}

package clockface

import (
	"fmt"
	"io"
	"log"
	"math"
	"time"
)

type Point struct {
	X float64
	Y float64
}

const (
	secondHandLength = float64(90)
	minuteHandLength = float64(80)
	hourHandLength   = float64(50)
	clockCenterX     = float64(150)
	clockCenterY     = float64(150)
)

const (
	bezel    = `<circle cx="150" cy="150" r="100" style="fill:#fff;stroke:#000;stroke-width:5px;"/>`
	svgEnd   = `</svg>`
	svgStart = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
	width="100%"
	height="100%"
	viewBox="0 0 300 300"
	version="2.0">`
)

func SecondHand(t time.Time) (point Point) {
	secondHandPoint := secondHandPoint(t)
	point.X = 150 + (secondHandPoint.X * secondHandLength)
	point.Y = 150 - (secondHandPoint.Y * secondHandLength)
	return
}

func secondsInRadians(t time.Time) float64 {
	return math.Pi / (30 / float64(t.Second()))
}

func minutesInRadians(t time.Time) float64 {
	return math.Pi / (30 / float64(t.Minute()))
}

func hoursInRadians(t time.Time) float64 {
	return math.Pi / (6 / (float64(t.Hour() % 12)))
}

func radianToHandPoint(radian float64) Point {
	return Point{
		X: math.Sin(radian),
		Y: math.Cos(radian),
	}
}

func secondHandPoint(t time.Time) Point {
	return radianToHandPoint(secondsInRadians(t))
}

func minuteHandPoint(t time.Time) Point {
	return radianToHandPoint(minutesInRadians(t))
}

func hourHandPoint(t time.Time) Point {
	return radianToHandPoint(hoursInRadians(t))
}

func makeHand(p Point, length float64) Point {
	p = Point{p.X * length, p.Y * -length}
	return Point{p.X + clockCenterX, p.Y + clockCenterY}
}

func secondHand(w io.Writer, t time.Time) {
	p := makeHand(secondHandPoint(t), secondHandLength)
	_, err := fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#f00;stroke-width:3px;"/>`, p.X, p.Y)
	if err != nil {
		log.Panicf("error printing to the writer: %s", err)
	}
}

func minuteHand(w io.Writer, t time.Time) {
	p := makeHand(minuteHandPoint(t), minuteHandLength)
	_, err := fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#000;stroke-width:3px;"/>`, p.X, p.Y)
	if err != nil {
		log.Panicf("error printing to the writer: %s", err)
	}
}

func hourHand(w io.Writer, t time.Time) {
	p := makeHand(hourHandPoint(t), hourHandLength)
	_, err := fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#000;stroke-width:3px;"/>`, p.X, p.Y)
	if err != nil {
		log.Panicf("error printing to the writer: %s", err)
	}
}

func SVGWriter(w io.Writer, t time.Time) {
	_, err := io.WriteString(w, svgStart)
	if err != nil {
		log.Panicf("error printing to the writer: %s", err)
	}
	_, err = io.WriteString(w, bezel)
	if err != nil {
		log.Panicf("error printing to the writer: %s", err)
	}
	secondHand(w, t)
	minuteHand(w, t)
	hourHand(w, t)
	_, err = io.WriteString(w, svgEnd)
	if err != nil {
		log.Panicf("error printing to the writer: %s", err)
	}
}

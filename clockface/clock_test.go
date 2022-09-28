package clockface

import (
	"bytes"
	"encoding/xml"
	"math"
	"reflect"
	"testing"
	"time"
)

type Line struct {
	X1 float64 `xml:"x1,attr"`
	Y1 float64 `xml:"y1,attr"`
	X2 float64 `xml:"x2,attr"`
	Y2 float64 `xml:"y2,attr"`
}

type Circle struct {
	Cx float64 `xml:"cx,attr"`
	Cy float64 `xml:"cy,attr"`
	R  float64 `xml:"r,attr"`
}

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	Version string   `xml:"version,attr"`
	Circle  Circle   `xml:"circle"`
	Line    []Line   `xml:"line"`
}

func simpleTime(hours, minutes, seconds int) time.Time {
	return time.Date(312, time.October, 28, hours, minutes, seconds, 0, time.UTC)
}

func testName(t time.Time) string {
	return t.Format("15:04:05")
}

func roughlyEqualFloat64(a, b float64) bool {
	const equalityThreshold = 1e-7
	return math.Abs(a-b) < equalityThreshold
}

func roughlyEqualPoint(a, b Point) bool {
	return roughlyEqualFloat64(a.X, b.X) &&
		roughlyEqualFloat64(a.Y, b.Y)
}

func TestSecondHand(t *testing.T) {

	var testCase = []struct {
		Time time.Time
		Want Point
	}{
		{
			Time: simpleTime(0, 0, 0),
			Want: Point{X: 150, Y: 150 - 90},
		},
		{
			Time: simpleTime(0, 0, 15),
			Want: Point{X: 150 + 90, Y: 150},
		},
		{
			Time: simpleTime(0, 0, 30),
			Want: Point{X: 150, Y: 150 + 90},
		},
		{
			Time: simpleTime(0, 0, 45),
			Want: Point{X: 150 - 90, Y: 150},
		},
	}

	for _, test := range testCase {
		t.Run(testName(test.Time), func(t *testing.T) {
			got := SecondHand(test.Time)
			if !roughlyEqualPoint(got, test.Want) {
				t.Errorf("Got %v, wanted %v", got, test.Want)
			}
		})
	}
}

func TestSecondsInRadians(t *testing.T) {
	cases := []struct {
		time  time.Time
		angle float64
	}{
		{simpleTime(0, 0, 0), 0},
		{simpleTime(0, 0, 7), (math.Pi / 30) * 7},
		{simpleTime(0, 0, 30), math.Pi},
		{simpleTime(0, 0, 45), (math.Pi / 2) * 3},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondsInRadians(c.time)
			if !roughlyEqualFloat64(got, c.angle) {
				t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
			}
		})
	}
}

func TestSecondHandVector(t *testing.T) {
	cases := []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(0, 0, 15), Point{1, 0}},
		{simpleTime(0, 0, 30), Point{0, -1}},
		{simpleTime(0, 0, 45), Point{-1, 0}},
		{simpleTime(0, 0, 60), Point{0, 1}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondHandPoint(c.time)
			if !roughlyEqualPoint(got, c.point) {
				t.Fatalf("Wanted %v Point, but got %v", c.point, got)
			}
		})
	}
}

func containsLine(l Line, ls []Line) bool {
	for _, line := range ls {
		if reflect.DeepEqual(l, line) {
			return true
		}
	}
	return false
}

func TestSVGWriter(t *testing.T) {

	cases := []struct {
		time time.Time
		line Line
	}{
		{
			time: simpleTime(0, 0, 0),
			line: Line{clockCenterX, clockCenterY, clockCenterX, clockCenterY - secondHandLength},
		},
		{
			time: simpleTime(0, 0, 15),
			line: Line{clockCenterX, clockCenterY, clockCenterX + secondHandLength, clockCenterY},
		},
		{
			time: simpleTime(0, 0, 30),
			line: Line{clockCenterX, clockCenterY, clockCenterX, clockCenterY + secondHandLength},
		},
		{
			time: simpleTime(0, 0, 45),
			line: Line{clockCenterX, clockCenterY, clockCenterX - secondHandLength, clockCenterY},
		},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			b := bytes.Buffer{}
			SVGWriter(&b, c.time)

			svg := SVG{}
			err := xml.Unmarshal(b.Bytes(), &svg)
			if err != nil {
				panic("error unmarshalling xml file: " + err.Error())
			}

			if !containsLine(c.line, svg.Line) {
				t.Errorf("expected to find the second hand with x2 of %+v and y2 of %+v, in the SVG output %v", c.line.X2, c.line.Y2, b.String())
			}
		})
	}
}

func TestSVGWriterMinuteHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{
			time: simpleTime(0, 0, 0),
			line: Line{clockCenterX, clockCenterY, clockCenterX, clockCenterY - minuteHandLength},
		},
		{
			time: simpleTime(0, 15, 0),
			line: Line{clockCenterX, clockCenterY, clockCenterX + minuteHandLength, clockCenterY},
		},
		{
			time: simpleTime(0, 30, 0),
			line: Line{clockCenterX, clockCenterY, clockCenterX, clockCenterY + minuteHandLength},
		},
		{
			time: simpleTime(0, 45, 0),
			line: Line{clockCenterX, clockCenterY, clockCenterX - minuteHandLength, clockCenterY},
		},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			b := bytes.Buffer{}
			SVGWriter(&b, c.time)

			svg := SVG{}
			err := xml.Unmarshal(b.Bytes(), &svg)
			if err != nil {
				panic("error unmarshalling xml file: " + err.Error())
			}

			if !containsLine(c.line, svg.Line) {
				t.Errorf("Expected to find the minute hand in line %+v, in the SVG lines %+v", c.line, svg.Line)
			}
		})
	}
}

func TestSVGWriterHourHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{
			time: simpleTime(0, 0, 0),
			line: Line{clockCenterX, clockCenterY, clockCenterX, clockCenterY - hourHandLength},
		},
		{
			time: simpleTime(3, 0, 0),
			line: Line{clockCenterX, clockCenterY, clockCenterX + hourHandLength, clockCenterY},
		},
		{
			time: simpleTime(6, 0, 0),
			line: Line{clockCenterX, clockCenterY, clockCenterX, clockCenterY + hourHandLength},
		},
		{
			time: simpleTime(9, 0, 0),
			line: Line{clockCenterX, clockCenterY, clockCenterX - hourHandLength, clockCenterY},
		},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			b := bytes.Buffer{}
			SVGWriter(&b, c.time)

			svg := SVG{}
			err := xml.Unmarshal(b.Bytes(), &svg)
			if err != nil {
				panic("error unmarshalling xml file: " + err.Error())
			}

			if !containsLine(c.line, svg.Line) {
				t.Errorf("Expected to find the minute hand in line %+v, in the SVG lines %+v", c.line, svg.Line)
			}
		})
	}
}

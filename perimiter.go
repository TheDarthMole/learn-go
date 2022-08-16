package main

import (
	"math"
)

// Rectangle methods

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

type Triangle struct {
	Height float64
	Base   float64
}

type Shape interface {
	Area() float64
}

func (rect Rectangle) Perimeter() float64 {
	return 2 * (rect.Width + rect.Height)
}

func (rect Rectangle) Area() float64 {
	return rect.Width * rect.Height
}

// Circle methods

func (circle Circle) Area() float64 {
	return math.Pi * math.Pow(circle.Radius, 2)
}

// Triangle methods

func (triangle Triangle) Area() float64 {
	return triangle.Base * triangle.Height * 0.5
}

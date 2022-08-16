package shapes

import (
	"testing"
)

func TestPerimeter(t *testing.T) {
	rect := Rectangle{Width: 32, Height: 64}
	got := rect.Perimeter()
	want := 192.0

	if got != want {
		t.Errorf("Got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {

	// This contains the test cases and their expected values
	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle", shape: Rectangle{12, 6}, hasArea: 72},
		{name: "Circle", shape: Circle{10}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{12, 6}, hasArea: 36},
	}

	t.Run("Table driven tests for Rectangle Circle and Triangle", func(t *testing.T) {
		for _, tt := range areaTests {

			t.Run(tt.name, func(t *testing.T) {
				got := tt.shape.Area()
				if got != tt.hasArea {
					t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
				}
			})

		}
	})

}

func BenchmarkPerimeter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rect := Rectangle{Width: 32, Height: 512}
		rect.Perimeter()
	}
}

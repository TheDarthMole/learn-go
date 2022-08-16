package iteration

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {

	t.Run("Test sum 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("Got %d want %d given %v", got, want, numbers)
		}
	})

}

func TestSumAll(t *testing.T) {
	t.Run("Sum 2 lists of integers", func(t *testing.T) {
		numbers := [][]int{{1, 2, 3, 4, 5}, {9, 8, 7, 6}}
		want := []int{15, 30}
		got := SumAll(numbers)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("Got %v want %v given %v", got, want, numbers)
		}
	})
}

func BenchmarkSumAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numbers := [][]int{{1, 2, 3, 4, 5, 6}, {9, 8, 7, 6, 5}}
		SumAll(numbers)
	}
}

func TestSumAllTails(t *testing.T) {

	t.Run("Make the sums of some slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2, 5}, []int{0, 9, 0}, []int{0, 0})
		want := []int{7, 9, 0}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got %v want %v", got, want)
		}
	})

	t.Run("Safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

package iteration

import "testing"

func TestSum(t *testing.T) {

	t.Run("Test sum 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("Got %d want %d given %v", got, want, numbers)
		}
	})

	//t.Run("Test sum 3 numbers", func(t *testing.T) {
	//	numbers := []int{5, 4, 3}
	//
	//	got := Sum(numbers)
	//	want := 12
	//
	//	if got != want {
	//		t.Errorf("Got %d want %d given %v", got, want, numbers)
	//	}
	//})

}

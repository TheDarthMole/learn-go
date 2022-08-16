package iteration

func Sum(numbers []int) (sum int) {
	for _, num := range numbers {
		sum += num
	}
	return
}

func SumAll(allLists [][]int) (sums []int) {
	lenList := len(allLists)
	sums = make([]int, lenList)
	for i, list := range allLists {
		sums[i] = Sum(list)
	}
	return
}

func SumAllTails(values ...[]int) (sums []int) {
	sums = make([]int, len(values))

	for index, numbers := range values {
		if len(numbers) == 0 {
			sums[index] = 0
		} else {
			tail := numbers[1:]
			sums[index] = Sum(tail)
		}
	}
	return
}

package iteration

func Sum(numbers []int) (sum int) {
	for _, num := range numbers {
		sum += num
	}
	return
}

//func SumAll(allLists [][]int) (retVal int[]){
//	return
//}

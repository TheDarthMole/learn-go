package iteration

func Repeat(character string, num int) (repeated string) {
	for i := 0; i < num; i++ {
		repeated += character
	}
	return
}

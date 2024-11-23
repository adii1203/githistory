package utils

func GetRange(start int, end int) []int {
	var r []int
	for i := start; i <= end; i++ {
		r = append(r, i)
	}
	return r
}

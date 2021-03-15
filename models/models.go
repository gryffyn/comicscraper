package models

func GenArray(first, last int) []int {
	a := make([]int, last-first+1)
	for i := range a {
		a[i] = first + i
	}
	return a
}

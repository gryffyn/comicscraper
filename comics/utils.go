package comics

import (
	"time"
)

func GenIntArray(first, last int) []int {
	a := make([]int, last-first+1)
	for i := range a {
		a[i] = first + i
	}
	return a
}

// Generates a slice of days.
func GenDateArray(first, last time.Time) []time.Time {
	var urls []time.Time
	for f := first; last.After(f); f = f.Add(time.Hour * 24) {
		urls = append(urls, f)
	}
	urls = append(urls, last)
	return urls
}

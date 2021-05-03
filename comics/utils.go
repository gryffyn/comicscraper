package comics

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

func GenIntArray(first, last int) []int {
	a := make([]int, last-first+1)
	for i := range a {
		a[i] = first + i
	}
	return a
}

// GenDateArray Generates a slice of days.
func GenDateArray(first, last time.Time) []time.Time {
	var dates []time.Time
	for f := first; last.After(f); f = f.Add(time.Hour * 24) {
		dates = append(dates, f)
	}
	dates = append(dates, last)
	return dates
}

func isPNG(input io.Reader) (io.Reader, bool, error) {
	buf := [4]byte{}

	n, err := io.ReadAtLeast(input, buf[:], len(buf))
	if err != nil {
		return nil, false, err
	}

	isGzip := buf[0] == 137 && buf[1] == 80 && buf[2] == 78 && buf[3] == 71
	return io.MultiReader(bytes.NewReader(buf[:n]), input), isGzip, nil
}

func Get(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	return resp, err
}

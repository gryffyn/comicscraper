package comics

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func GenQCURLs(first, last int) []string {
	var urls []string
	for f := first; f <= last; f++ {
		urls = append(urls, "https://questionablecontent.net/comics/"+strconv.Itoa(f)+".png")
	}
	return urls
}

// Now with wait groups!
func GetQCStripAll(arr []int, filepath string, bar *progressbar.ProgressBar) error {
	total := len(arr)
	size := total / 4
	rmdr := total % 4
	wg := &sync.WaitGroup{}
	var err error
	for i := 0; i < 4; i++ {
		wg.Add(1)
		start := i * size
		end := (i + 1) * size
		if i == 3 {
			end += rmdr
		}

		go func(start, end, i int) {
			for f := start; f < end; f++ {
				err = GetQCStrip(arr[f], filepath, bar)
			}
			wg.Done()
		}(start, end, i)
	}
	wg.Wait()
	return err
}

func GetQCStrip(num int, filepath string, bar *progressbar.ProgressBar) error {
	url := "https://questionablecontent.net/comics/" + strconv.Itoa(num) + ".png"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		url = "https://questionablecontent.net/comics/" + strconv.Itoa(num) + ".jpg"
		resp, err = http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	out, err := os.Create(filepath + fmt.Sprintf("%04d", num) + ".png")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	// log.Println("Downloaded file: " + filepath + strconv.Itoa(num) + ".png")
	bar.Add(1)
	return err
}

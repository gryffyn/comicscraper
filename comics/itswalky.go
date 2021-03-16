package comics

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Downloads all urls in arr into filepath. Now with wait groups!
func GetIWStripAll(arr []time.Time, filepath string, bar *progressbar.ProgressBar) error {
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
				err = GetIWStrip(arr[f], filepath, bar)
			}
			wg.Done()
		}(start, end, i)
	}
	wg.Wait()
	return err
}

// Downloads a single strip. Tries .gif first.
func GetIWStrip(strip time.Time, filepath string, bar *progressbar.ProgressBar) error {
	PNG := false
	layout := "2006-01-02"
	url := "https://www.itswalky.com/wp-content/uploads/" + strconv.Itoa(strip.Year()) + "/" + fmt.Sprintf("%02d", strip.Month()) + "/" + strip.Format(layout)
	resp, err := http.Get(url + ".gif")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		resp, err = http.Get(url + ".png")
		if err != nil {
			return err
		}
		PNG = true
		defer resp.Body.Close()
	}

	var out *os.File
	if PNG {
		out, err = os.Create(filepath + strip.Format(layout) + ".png")
	} else {
		out, err = os.Create(filepath + strip.Format(layout) + ".gif")
	}

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	// log.Println("Downloaded file: " + filepath + strconv.Itoa(num) + ".png")
	bar.Add(1)
	return err
}

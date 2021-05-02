package dlstrip

import (
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type dlStripDate func(time.Time, string, *progressbar.ProgressBar) error
type dlStripInt func(int, string, *progressbar.ProgressBar) error

func GetAllDate(arr []time.Time, filepath string, bar *progressbar.ProgressBar, dlstrip dlStripDate) error {
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
				err = dlstrip(arr[f], filepath, bar)
			}
			wg.Done()
		}(start, end, i)
	}
	wg.Wait()
	return err
}

func GetAllInt(arr []int, filepath string, bar *progressbar.ProgressBar, dlstrip dlStripInt) error {
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
				err = dlstrip(arr[f], filepath, bar)
			}
			wg.Done()
		}(start, end, i)
	}
	wg.Wait()
	return err
}

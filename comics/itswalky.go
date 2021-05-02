package comics

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
)

// GetIWStrip Downloads a single strip. Tries .gif first.
func GetIWStrip(strip time.Time, filepath string, bar *progressbar.ProgressBar) error {
	layout := "2006-01-02"
	url := "https://www.itswalky.com/wp-content/uploads/" + strconv.Itoa(strip.Year()) + "/" + fmt.Sprintf("%02d",
		strip.Month()) + "/" + strip.Format(layout) + ".gif"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var timedate string
	if resp.StatusCode == 404 {
		if strip.Month() == 12 {
			timedate = strconv.Itoa(strip.Year()+1) + "/" + fmt.Sprintf("%02d", 1)
		} else {
			timedate = strconv.Itoa(strip.Year()) + "/" + fmt.Sprintf("%02d", strip.Month()+1)
		}
		url = "https://www.itswalky.com/wp-content/uploads/" + timedate + "/" + strip.Format(layout) + ".gif"
		resp, err = http.Get(url)
		if err != nil {
			return err
		}
		if resp.StatusCode == 404 {
			url = "https://www.itswalky.com/wp-content/uploads/" + timedate + "/" + strip.Format(layout) + ".png"
			resp, err = http.Get(url)
			if err != nil {
				return err
			}
		}
		defer resp.Body.Close()
	}

	var out *os.File

	out, err = os.Create(filepath + strip.Format(layout) + url[len(url)-4:])

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	// log.Println("Downloaded file: " + filepath + strconv.Itoa(num) + ".png")
	bar.Add(1)
	return err
}

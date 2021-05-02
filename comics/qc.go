package comics

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/schollz/progressbar/v3"
)

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

	out, err := os.Create(filepath + fmt.Sprintf("%04d", num) + url[len(url)-4:])
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	// log.Println("Downloaded file: " + filepath + strconv.Itoa(num) + ".png")
	bar.Add(1)
	return err
}

package comics

import (
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/schollz/progressbar/v3"
)

//go:embed ggar-titles
var ggarTitles string
var ggarArray = strings.Split(ggarTitles, "\n")

func GetGGARStrip(stripIndex int, filepath string, bar *progressbar.ProgressBar) error {
	strip := ggarArray[stripIndex]
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.gogetaroomie.com/comic/"+strip, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	stripRegex := regexp.MustCompile(`title=.*src="(https:\/\/.*).*"\sid=`)
	parsed := stripRegex.FindStringSubmatch(string(body))

	respStrip, err := http.Get(parsed[1])
	if err != nil {
		return err
	}
	defer respStrip.Body.Close()

	var out *os.File
	out, err = os.Create(filepath + fmt.Sprintf("%04d", stripIndex) + "-" + strip + parsed[1][len(parsed[1])-4:])
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, respStrip.Body)

	bar.Add(1)
	return err
}

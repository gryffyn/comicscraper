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

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html"
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

func GetGGARText(stripName string, fileindex int, filepath string) error {
	resp, err := http.Get("https://www.gogetaroomie.com/ggar-rerun/" + stripName)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	getComment := regexp.MustCompile(`cc-newsbody"\>(.*)You can find these two strips in the old archive`)
	commentnobsp := strings.Replace(getComment.FindStringSubmatch(string(body))[1], " ", " ", -1)
	commentnobsp2 := strings.Replace(commentnobsp, "&nbsp;", " ", -1)
	cnbspbr := strings.Replace(commentnobsp2, `<br>`, "\n", -1)
	comment := strip.StripTags(cnbspbr)

	var out *os.File
	out, err = os.Create(filepath + fmt.Sprintf("%04d", fileindex) + "-" + stripName + ".txt")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.WriteString(comment)
	if err != nil {
		return err
	}
	return err
}

func GetGGARTitles(series string) ([]string, error) {
	var titles []string
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.gogetaroomie.com/"+series+"/archive", nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	doc, _ := html.Parse(resp.Body)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "option" {
			if n.Attr[0].Val != "" {
				titles = append(titles, strings.TrimPrefix(n.Attr[0].Val, "ggar-rerun/"))
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return titles, err
}

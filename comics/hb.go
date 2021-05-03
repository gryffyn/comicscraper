package comics

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html"
)

func GetHBStrip(strip int, filepath string, bar *progressbar.ProgressBar) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.headlessbliss.com/comic/page-"+strconv.Itoa(strip+1), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var url string
	doc, _ := html.Parse(resp.Body)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			if len(n.Attr) == 3 {
				if n.Attr[2].Val == "cc-comic" {
					url = n.Attr[1].Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	reqStrip, err := http.NewRequest("GET", url, nil)
	respStrip, err := client.Do(reqStrip)
	if err != nil {
		return err
	}
	defer respStrip.Body.Close()

	var out *os.File
	out, err = os.Create(filepath + fmt.Sprintf("%04d", strip+1) + url[len(url)-4:])
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, respStrip.Body)

	bar.Add(1)
	return err
}

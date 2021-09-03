package comics

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html"
)

func GetXKCDStrip(num int, filepath string, bar *progressbar.ProgressBar) error {
	url := "https://xkcd.com/" + strconv.Itoa(num) + "/"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	var imgUrl,alttext string

	doc, _ := html.Parse(resp.Body)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			if n.Attr[0].Key == "href" {
				if strings.HasPrefix(n.Attr[0].Val, "https://imgs.xkcd.com/comics/") {
					imgUrl = n.Attr[0].Val
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "img" {
			if len(n.Attr) > 1 {
				if n.Attr[1].Key == "title" {
					alttext = n.Attr[1].Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	resp.Body.Close()

	time.Sleep(1 * time.Second) // rate limit prevention

	if imgUrl != "" {
		respImg, err := http.Get(imgUrl)

		out, err := os.Create(filepath + fmt.Sprintf("%04d", num) + imgUrl[len(imgUrl)-4:])
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, respImg.Body)

		err = respImg.Body.Close()
		if err != nil {
			log.Fatalln("err downloading url '" + imgUrl + "'")
		}

		outText, err := os.Create(filepath + fmt.Sprintf("%04d", num) + imgUrl[len(imgUrl)-4:] + ".alttext")
		if err != nil {
			return err
		}
		defer outText.Close()

		_,err = outText.WriteString(alttext)
		if err != nil {
			return err
		}

		bar.Add(1)
	} else {
		err = errors.New("could not download: imgUrl is '" + imgUrl + "'")
	}

	return err
}
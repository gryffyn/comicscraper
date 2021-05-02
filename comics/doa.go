package comics

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html"
)

const StartDate_DOA = "2010-09-06"

func GetDOAStrip(strip time.Time, filepath string, bar *progressbar.ProgressBar) error {
	layout := "2006-01-02"
	urls, _, _, err := GetDOAURLS(strip.Year())

	url := urls[strip.Format(layout)]
	if url == "DNE" {
		return nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		resp, err = http.Get(url[0:len(url)-6] + ".png")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	read, ispng, err := isPNG(resp.Body)
	if !ispng {
		return nil
	}

	var out *os.File
	filename := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	out, err = os.Create(filepath + filename)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, read)

	bar.Add(1)
	return err
}

func GetDOAURLS(year int) (map[string]string, []string, int64, error) {
	var comics []string
	vals := make(map[string]string)

	start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	layout := "2006-01-02"
	var keys []string

	url := "https://www.dumbingofage.com/archive/?archive_year=" + strconv.Itoa(year)
	resp, err := http.Get(url)

	defer resp.Body.Close()

	if year <= 2013 {
		z := html.NewTokenizer(resp.Body)
		for {
			tt := z.Next()
			if tt == html.ErrorToken {
				break
			}
			_, val, _ := z.TagAttr()
			// fmt.Println("Key:", string(key), "Val:", string(val))
			if string(val) == "cpcal-day" {
				z.Next()
				z.Next()
				_, url, _ := z.TagAttr()
				comics = append(comics, string(url))
			}
		}
	} else {
		doc, _ := html.Parse(resp.Body)
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				if n.Attr[0].Key == "href" {
					if strings.HasPrefix(n.Attr[0].Val, "https://www.dumbingofage.com/"+strconv.Itoa(year)+"/comic/") {
						comics = append(comics, n.Attr[0].Val)
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
	}

	var days int64 = 0
	for _, f := range comics {
		var url string
		if f == "" {
			url = "DNE"
		} else {
			days++
			sr := strings.Split(f, "/")
			url = "https://www.dumbingofage.com/comics/" + start.Format(layout) + "-" + sr[len(sr)-2] + ".png"
		}
		keys = append(keys, start.Format(layout))
		vals[start.Format(layout)] = url
		start = start.Add(time.Hour * 24)
	}

	return vals, keys, days, err
}

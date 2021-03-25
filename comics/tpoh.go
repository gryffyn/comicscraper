package comics

import (
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func getChapters() ([]string, error) {
	var chapters []string
	chapters = append(chapters, "The Hook")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://jolleycomics.com/TPoH/The%20Hook/", nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	doc, _ := html.Parse(resp.Body)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "option" {
			if len(n.Attr) == 2 {
				ch := strings.TrimSuffix(strings.TrimPrefix(n.Attr[1].Val, "/TPoH/"), "/")
				chapters = append(chapters, ch)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return chapters, err
}

func getComicURL(number int, chapter string) (string, error) {
	client := &http.Client{}
	var url string
	req, err := http.NewRequest("GET", "http://jolleycomics.com/TPoH/"+chapter+"/"+strconv.Itoa(number)+"/", nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	doc, _ := html.Parse(resp.Body)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			c := n.Attr[0].Val
			if strings.Contains(c, "/comics/") {
				url = "http://jolleycomics.com" + c
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return url, err
}

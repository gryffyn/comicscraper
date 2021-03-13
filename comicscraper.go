package comicscraper

import (
	"bufio"
	"fmt"
	"go/scanner"
	"io"
	"net/http"
	"os"
)

func main() {
	fileUrl := "https://golangcode.com/logo.svg"
	err := DownloadFile("logo.svg", fileUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded: " + fileUrl)
}

// Downloads file to current dir.
func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

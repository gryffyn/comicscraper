package models

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GetQCStrip(num int, filepath string) error {
	url := "https://questionablecontent.net/comics/" + strconv.Itoa(num) + ".png"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath + strconv.Itoa(num) + ".png")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	log.Println("Downloaded file: " + filepath + strconv.Itoa(num) + ".png")
	return err
}

func GetQCStripAll(first, last int, filepath string) error {
	var err error
	for f := first; f <= last; f++ {
		err = GetQCStrip(f, filepath)
	}
	return err
}

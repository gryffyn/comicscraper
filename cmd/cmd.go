package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"git.neveris.one/gryffyn/comicscraper/comics"
	"git.neveris.one/gryffyn/comicscraper/dlstrip"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

func Run() {
	var comic, dir, first, last string
	var text bool

	app := &cli.App{
		Name:      "comicscraper",
		Version:   "v0.0.1-alpha",
		Compiled:  time.Now(),
		Copyright: "(c) 2021 gryffyn",
		Usage:     "download comic images. Date format is 'YYYY-MM-DD'.",
		UsageText: "comicscraper [arguments]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "comic",
				Aliases:     []string{"c"},
				Usage:       "name of the comic to download",
				Destination: &comic,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "directory",
				Aliases:     []string{"d"},
				Value:       ".",
				Usage:       "directory to download into",
				Destination: &dir,
			},
			&cli.StringFlag{
				Name:        "first",
				Aliases:     []string{"f"},
				Usage:       "number/date of the comic, or first if downloading multiple",
				Destination: &first,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "last",
				Aliases:     []string{"l"},
				Usage:       "number/date of the last comic",
				Destination: &last,
			},
			&cli.BoolFlag{
				Name:        "text",
				Aliases:     []string{"t"},
				Usage:       "associated text of comic",
				Destination: &text,
			},
		},
		Action: func(c *cli.Context) error {
			var err error
			dir = fixPath(dir)
			if strings.ToLower(c.String("comic")) == "qc" {
				fi, _ := strconv.Atoi(first)
				if last == "" {
					bar := progressbar.Default(1)
					err = comics.GetQCStrip(fi, dir, bar)
				} else {
					li, _ := strconv.Atoi(last)
					max := li - fi + 1
					bar := progressbar.Default(int64(max))
					err = dlstrip.GetAllInt(comics.GenIntArray(fi, li), dir, bar, comics.GetQCStrip)
				}
			} else if strings.ToLower(c.String("comic")) == "iw" {
				strips, firstDate, days, _ := date(first, last)
				if last == "" {
					bar := progressbar.Default(1)
					err = comics.GetIWStrip(firstDate, dir, bar)
				} else {
					bar := progressbar.Default(days + 1)
					err = dlstrip.GetAllDate(strips, dir, bar, comics.GetIWStrip)
				}
			} else if strings.ToLower(c.String("comic")) == "doa" {
				strips, firstDate, _, _ := date(first, last)
				if last == "" {
					bar := progressbar.Default(1)
					err = comics.GetDOAStrip(firstDate, dir, bar)
				} else {
					bar := progressbar.Default(-1)
					err = dlstrip.GetAllDate(strips, dir, bar, comics.GetDOAStrip)
				}
			} else if strings.ToLower(c.String("comic")) == "ggar" {
				fi, _ := strconv.Atoi(first)
				if last == "" {
					bar := progressbar.Default(1)
					err = comics.GetGGARStrip(fi, dir, bar)
				} else {
					li, _ := strconv.Atoi(last)
					max := li - fi + 1
					bar := progressbar.Default(int64(max - 1))
					err = dlstrip.GetAllInt(comics.GenIntArray(fi, li), dir, bar, comics.GetGGARStrip)
				}
			} else if strings.ToLower(c.String("comic")) == "hb" {
				fi, _ := strconv.Atoi(first)
				if last == "" {
					bar := progressbar.Default(1)
					err = comics.GetHBStrip(fi, dir, bar)
				} else {
					li, _ := strconv.Atoi(last)
					max := li - fi + 1
					bar := progressbar.Default(int64(max - 1))
					err = dlstrip.GetAllInt(comics.GenIntArray(fi, li), dir, bar, comics.GetHBStrip)
				}
			} else {
				log.Fatalln("Comic not found.")
			}
			fmt.Println("\nFinished downloading.")
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func fixPath(path string) string {
	if runtime.GOOS == "windows" {
		if path[len(path)-1] != '\\' {
			path += "\\"
		}
	} else {
		if path[len(path)-1] != '/' {
			path += "/"
		}
	}
	return path
}

func checkDate(first, last string) bool {
	layout := "2006-01-02"
	firstDate, _ := time.Parse(layout, first)
	lastDate, _ := time.Parse(layout, last)
	return firstDate.Before(lastDate)
}

func date(first, last string) ([]time.Time, time.Time, int64, error) {
	if first == "2014-02-29" { // damn you, willis
		first = "2014-02-28"
	} else if last == "2014-02-29" {
		last = "2014-02-28"
	}

	var err error
	if last != "" {
		if !checkDate(first, last) {
			err = errors.New("first date is not before last date")
		}
	}
	layout := "2006-01-02"
	firstDate, err := time.Parse(layout, first)
	lastDate, err := time.Parse(layout, last)
	if err != nil {
		log.Fatalln(err)
	}
	days := int64(lastDate.Sub(firstDate).Hours() / 24)
	strips := comics.GenDateArray(firstDate, lastDate)

	return strips, firstDate, days, err
}

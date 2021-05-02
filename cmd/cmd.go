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
		},
		Action: func(c *cli.Context) error {
			var err error
			dir = fixPath(dir)

			if first == "2014-02-29" { // damn you, willis
				first = "2014-02-28"
			} else if last == "2014-02-29" {
				last = "2014-02-28"
			}

			if last != "" {
				if !checkDate(first, last) {
					return errors.New("first date is not before last date")
				}
			}
			layout := "2006-01-02"
			firstDate, err := time.Parse(layout, first)
			lastDate, err := time.Parse(layout, last)
			if err != nil {
				log.Fatalln(err)
			}
			days := lastDate.Sub(firstDate).Hours() / 24
			strips := comics.GenDateArray(firstDate, lastDate)
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
				if last == "" {
					bar := progressbar.Default(1)
					err = comics.GetIWStrip(firstDate, dir, bar)
				} else {
					bar := progressbar.Default(int64(days + 1))
					err = dlstrip.GetAllDate(strips, dir, bar, comics.GetIWStrip)
				}
			} else if strings.ToLower(c.String("comic")) == "doa" {
				if last == "" {
					bar := progressbar.Default(1)
					err = comics.GetDOAStrip(firstDate, dir, bar)
				} else {
					bar := progressbar.Default(-1)
					err = dlstrip.GetAllDate(strips, dir, bar, comics.GetDOAStrip)
				}
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

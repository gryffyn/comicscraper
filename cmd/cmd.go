package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"git.neveris.one/gryffyn/comicscraper/comics"
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
			if !checkDate(first, last) {
				log.Fatal("First date is not before the last date.")
			}
			if strings.ToLower(c.String("comic")) == "qc" {
				li, _ := strconv.Atoi(last)
				fi, _ := strconv.Atoi(first)
				max := li - fi + 1
				bar := progressbar.Default(int64(max))
				if li == 0 {
					err = comics.GetQCStrip(fi, dir, bar)
				} else {
					err = comics.GetQCStripAll(comics.GenIntArray(fi, li), dir, bar)
				}
			} else if strings.ToLower(c.String("comic")) == "iw" {
				layout := "2006-01-02"
				firstDate, _ := time.Parse(layout, first)
				lastDate, _ := time.Parse(layout, last)
				days := lastDate.Sub(firstDate).Hours() / 24
				bar := progressbar.Default(int64(days + 1))
				if last == "" {
					err = comics.GetIWStrip(firstDate, dir, bar)
				} else {
					strips := comics.GenDateArray(firstDate, lastDate)
					err = comics.GetIWStripAll(strips, dir, bar)
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

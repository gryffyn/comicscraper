package cmd

import (
	"log"
	"os"
	"strings"

	"git.neveris.one/gryffyn/comicscraper/models"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

func Run() {
	var comic, dir string
	var first, last int

	app := &cli.App{
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
				Usage:       "directory to download into (default: CWD)",
				Destination: &dir,
			},
			&cli.IntFlag{
				Name:        "first",
				Aliases:     []string{"f"},
				Usage:       "number of the comic, or first if downloading multiple",
				Destination: &first,
				Required:    true,
			},
			&cli.IntFlag{
				Name:        "last",
				Aliases:     []string{"l"},
				Usage:       "number of the last comic",
				Destination: &last,
			},
		},
		Action: func(c *cli.Context) error {
			var err error
			if strings.ToLower(c.String("comic")) == "qc" {
				max := last - first + 1
				bar := progressbar.Default(int64(max))
				if last == 0 {
					err = models.GetQCStrip(first, dir, bar)
				} else {
					err = models.GetQCStripAll(models.GenArray(first, last), dir, bar)
				}
			}
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

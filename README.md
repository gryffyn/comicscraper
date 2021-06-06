# comicscraper
[![Build Status](https://ci.neveris.one/api/badges/gryffyn/comicscraper/status.svg?ref=refs/heads/main)](https://ci.neveris.one/gryffyn/comicscraper)
[![Go Report Card](https://goreportcard.com/badge/git.neveris.one/gryffyn/comicscraper)](https://goreportcard.com/report/git.neveris.one/gryffyn/comicscraper)
## Installing

`go get git.neveris.one/gryffyn/comicscraper`

## Supported comics

- Dumbing of Age
- Go Get a Roomie
- Headless Bliss
- It's Walky!
- Questionable Content
- The Property of Hate

*working on a system so that this base application won't download anything by itself*

## Usage

```
§ comicscraper -h
NAME:
   comicscraper - download comic images. Date format is 'YYYY-MM-DD'.

USAGE:
   comicscraper [arguments]

VERSION:
   v0.0.1-alpha

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --comic value, -c value      name of the comic to download
   --directory value, -d value  directory to download into (default: ".")
   --first value, -f value      number/date of the comic, or first if downloading multiple
   --last value, -l value       number/date of the last comic
   --help, -h                   show help (default: false)
   --version, -v                print the version (default: false)

COPYRIGHT:
   (c) 2021 gryffyn
```

## License
See `LICENSE` for details.

The author takes no responsibility for the use of this software.
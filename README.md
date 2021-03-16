# comicscraper
[![Build Status](https://ci.neveris.one/api/badges/gryffyn/comicscraper/status.svg?ref=refs/heads/main)](https://ci.neveris.one/gryffyn/comicscraper)
## Installing

`go install git.neveris.one/gryffyn/comicscraper`

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
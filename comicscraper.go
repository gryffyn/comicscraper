package main

import (
	"fmt"

	"git.neveris.one/gryffyn/comicscraper/cmd"
)

func main() {
	cmd.Run()
	fmt.Print("\nFinished downloading.")
}

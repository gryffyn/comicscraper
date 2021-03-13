package main

import (
	"fmt"
	"log"

	"comicscraper/models"
)

func main() {
	fp := "_test/QC/"
	err := models.GetQCStripAll(1, 4322, fp)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Finished downloading")
}

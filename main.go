package main

import (
	"fmt"
	"log"

	"github.com/soloviev1d/cliautotestsuite/algorithm"
	"github.com/soloviev1d/cliautotestsuite/config"
)


func main() {
	config.FlagsInit()
	report, err := algorithm.RunTestSuites()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(report)
}

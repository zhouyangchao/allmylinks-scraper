package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zhouyangchao/allmylinks-scraper/allmylinks"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a username or full URL as an argument")
	}

	aml := allmylinks.NewAllMyLinks("")

	input := os.Args[1]
	var username, url string

	if strings.HasPrefix(input, "https://") {
		url = input
	} else {
		username = input
	}

	userInfo, err := aml.ScrapeUserInfo(username, url)
	if err != nil {
		log.Fatalf("Error scraping user info: %v", err)
	}

	fmt.Print(userInfo.String())
}

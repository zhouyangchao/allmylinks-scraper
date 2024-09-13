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

	input := os.Args[1]
	var username, url string

	if strings.HasPrefix(input, "https://") {
		url = input
	} else {
		username = input
	}

	userInfo, err := allmylinks.ScrapeUserInfo(username, url)
	if err != nil {
		log.Fatalf("Error scraping user info: %v", err)
	}

	fmt.Printf("Username: %s\n", userInfo.Username)
	fmt.Printf("Avatar URL: %s\n", userInfo.AvatarURL)
	fmt.Printf("Display Name: %s\n", userInfo.DisplayName)
	fmt.Printf("Birthday: %s\n", userInfo.Birthday)
	fmt.Printf("Bio: %s\n", userInfo.Bio)
	fmt.Printf("Content: %s\n", userInfo.Content)
	fmt.Printf("Location: %s\n", userInfo.Location)
	fmt.Printf("Profile Views: %s\n", userInfo.ProfileViews)
	fmt.Printf("Last Online: %s\n", userInfo.LastOnline.Format("2006-01-02 15:04:05"))
	fmt.Printf("QR Code URL: %s\n", userInfo.QRCodeURL)
	fmt.Printf("\nLinks:\n")
	for _, link := range userInfo.Links {
		fmt.Printf("- %s: %s (%s) (%s)\n", link.Title, link.URL, link.URLText, link.ConnectedStatus)
	}
}

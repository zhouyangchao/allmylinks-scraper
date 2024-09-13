package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/zhouyangchao/allmylinks-scraper/allmylinks"
)

func main() {
	// 检查是否提供了命令行参数
	if len(os.Args) < 2 {
		log.Fatal("Please provide a username or full URL as an argument")
	}

	input := os.Args[1]
	var username string

	// 判断输入是URL还是用户名
	if strings.HasPrefix(input, "https://") || strings.HasPrefix(input, "http://") {
		// 如果是URL，从中提取用户名
		parsedURL, err := url.Parse(input)
		if err != nil {
			log.Fatalf("Invalid URL: %v", err)
		}
		parts := strings.Split(parsedURL.Path, "/")
		if len(parts) > 1 {
			username = parts[len(parts)-1]
		} else {
			log.Fatal("Invalid URL format")
		}
	} else {
		// 如果不是URL，直接使用输入作为用户名
		username = input
	}

	userInfo, err := allmylinks.ScrapeUserInfo(username)
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

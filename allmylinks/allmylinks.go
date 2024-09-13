package allmylinks

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type AllMyLinks struct {
	ProxyURL string
}

type UserInfo struct {
	Username     string     `json:"Username" bson:"Username"`
	AvatarURL    string     `json:"AvatarURL" bson:"AvatarURL"`
	DisplayName  string     `json:"DisplayName" bson:"DisplayName"`
	Birthday     *string    `json:"Birthday,omitempty" bson:"Birthday,omitempty"`
	Bio          *string    `json:"Bio,omitempty" bson:"Bio,omitempty"`
	Content      *string    `json:"Content,omitempty" bson:"Content,omitempty"`
	Location     *string    `json:"Location,omitempty" bson:"Location,omitempty"`
	ProfileViews *string    `json:"ProfileViews,omitempty" bson:"ProfileViews,omitempty"`
	LastOnline   *time.Time `json:"LastOnline,omitempty" bson:"LastOnline,omitempty"`
	QRCodeURL    string     `json:"QRCodeURL" bson:"QRCodeURL"`
	Links        []Link     `json:"Links,omitempty" bson:"Links,omitempty"`
}

type Link struct {
	Title           string  `json:"Title" bson:"Title"`
	URL             string  `json:"URL" bson:"URL"`
	URLText         *string `json:"URLText,omitempty" bson:"URLText,omitempty"`
	ConnectedStatus *string `json:"ConnectedStatus,omitempty" bson:"ConnectedStatus,omitempty"`
}

func NewAllMyLinks(proxyURL string) *AllMyLinks {
	return &AllMyLinks{ProxyURL: proxyURL}
}

func (l *Link) String() string {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("- %s - %s", l.Title, l.URL))
	if l.URLText != nil {
		buffer.WriteString(fmt.Sprintf(" (%s)", *l.URLText))
	}
	if l.ConnectedStatus != nil {
		buffer.WriteString(fmt.Sprintf(" [%s]", *l.ConnectedStatus))
	}
	buffer.WriteString("\n")
	return buffer.String()
}

func (u *UserInfo) String() string {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("Username: %s\n", u.Username))
	buffer.WriteString(fmt.Sprintf("AvatarURL: %s\n", u.AvatarURL))
	buffer.WriteString(fmt.Sprintf("DisplayName: %s\n", u.DisplayName))
	if u.Birthday != nil {
		buffer.WriteString(fmt.Sprintf("Birthday: %s\n", *u.Birthday))
	}
	if u.Bio != nil {
		buffer.WriteString(fmt.Sprintf("Bio: %s\n", *u.Bio))
	}
	if u.Content != nil {
		buffer.WriteString(fmt.Sprintf("Content: %s\n", *u.Content))
	}
	if u.Location != nil {
		buffer.WriteString(fmt.Sprintf("Location: %s\n", *u.Location))
	}
	if u.ProfileViews != nil {
		buffer.WriteString(fmt.Sprintf("ProfileViews: %s\n", *u.ProfileViews))
	}
	if u.LastOnline != nil {
		buffer.WriteString(fmt.Sprintf("LastOnline: %s\n", u.LastOnline.Format("2006-01-02 15:04:05")))
	}
	buffer.WriteString(fmt.Sprintf("QRCodeURL: %s\n", u.QRCodeURL))
	buffer.WriteString("\nLinks:\n")
	for _, link := range u.Links {
		buffer.WriteString(link.String())
	}
	return buffer.String()
}

func (a *AllMyLinks) ScrapeUserInfo(username string, url string) (*UserInfo, error) {
	if username != "" {
		url = fmt.Sprintf("https://allmylinks.com/%s", username)
	}

	body, err := a.fetchHTMLDocument(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	userInfo := &UserInfo{Username: username}
	parseHTML(doc, userInfo)

	err = a.getProfileViews(doc, userInfo)
	if err != nil {
		return nil, err
	}

	// Remove duplicate links
	userInfo.Links = removeDuplicates(userInfo.Links)

	return userInfo, nil
}

func (a *AllMyLinks) fetchHTMLDocument(_url string) (io.ReadCloser, error) {
	client := &http.Client{}

	if a.ProxyURL != "" {
		proxyURL, err := url.Parse(a.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func parseHTML(n *html.Node, userInfo *UserInfo) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "img":
			if hasAttrValue(n, "alt", "Profile avatar") {
				userInfo.AvatarURL = getAttr(n, "src")
			}
		case "span":
			if hasClass(n, "profile-username") {
				userInfo.DisplayName = getTextContent(n)
			} else if hasClass(n, "last_online") {
				timestamp, err := strconv.ParseInt(getAttr(n, "data-x-timestamp"), 10, 64)
				if err == nil {
					lastOnline := time.Unix(timestamp, 0)
					userInfo.LastOnline = &lastOnline
				}
			}
		case "div":
			if hasClass(n, "about-section__birthday") {
				birthday := getTextContent(n)
				userInfo.Birthday = &birthday
			} else if hasClass(n, "about-section__location") {
				location := getAttr(n, "title")
				userInfo.Location = &location
			} else if hasClass(n, "about-section__content") {
				content := getTextContent(n)
				userInfo.Content = &content
			} else if hasClass(n, "simple-text") {
				link := parseLink(n)
				if link != nil {
					userInfo.Links = append(userInfo.Links, *link)
				}
			}
		case "p":
			if hasClass(n, "profile-bio") {
				bio := getTextContent(n)
				userInfo.Bio = &bio
			}
		case "a":
			if hasClass(n, "btn", "btn-link", "btn-qr") {
				userInfo.QRCodeURL = getAttr(n, "href")
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseHTML(c, userInfo)
	}
}

func hasClass(n *html.Node, classes ...string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			nodeClasses := strings.Fields(attr.Val)
			for _, class := range classes {
				found := false
				for _, nodeClass := range nodeClasses {
					if nodeClass == class {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
			return true
		}
	}
	return false
}

func hasAttrValue(n *html.Node, key, value string) bool {
	for _, attr := range n.Attr {
		if attr.Key == key && attr.Val == value {
			return true
		}
	}
	return false
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += getTextContent(c)
	}
	return strings.TrimSpace(result)
}

func removeDuplicates(links []Link) []Link {
	seen := make(map[string]bool)
	result := []Link{}
	for _, link := range links {
		if !seen[link.URL] {
			seen[link.URL] = true
			result = append(result, link)
		}
	}
	return result
}

func parseLink(n *html.Node) *Link {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			url := getAttr(c, "data-x-url")
			if url == "" {
				continue // Skip links without a valid URL
			}
			if url == "javascript: void(0);" {
				url = "N/A"
			}
			var title, text, connectedStatus string
			for gc := c.FirstChild; gc != nil; gc = gc.NextSibling {
				if gc.Type == html.ElementNode && gc.Data == "span" {
					if hasClass(gc, "link-title") {
						title = getTextContent(gc)
					} else if hasClass(gc, "link-text") {
						text = getTextContent(gc)
					} else if hasClass(gc, "connected-link-label") {
						connectedStatus = getTextContent(gc)
					}
				}
			}
			return &Link{Title: title, URL: url, URLText: &text, ConnectedStatus: &connectedStatus}
		}
	}
	return nil
}

func (a *AllMyLinks) getProfileViews(doc *html.Node, userInfo *UserInfo) error {
	profileViewsURL := findProfileViewsURL(doc)
	if profileViewsURL == "" {
		return nil
	}

	body, err := a.fetchHTMLDocument(profileViewsURL)
	if err != nil {
		return fmt.Errorf("failed to fetch profile views: %w", err)
	}
	defer body.Close()

	profileViewsBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read profile views: %w", err)
	}
	profileViews := strings.TrimSpace(string(profileViewsBytes))
	userInfo.ProfileViews = &profileViews

	return nil
}

func findProfileViewsURL(doc *html.Node) string {
	var profileViewsURL string
	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			content := getTextContent(n)
			if strings.Contains(content, "$.get(\"/profile/views?id=") {
				re := regexp.MustCompile(`/profile/views\?id=(\d+)`)
				matches := re.FindStringSubmatch(content)
				if len(matches) > 1 {
					profileViewsURL = "https://allmylinks.com" + matches[0]
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(doc)
	return profileViewsURL
}

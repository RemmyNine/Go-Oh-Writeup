package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	// List of RSS feed URLs to scrape
	urls := []string{
		"https://medium.com/feed/tag/bug-bounty-writeup",
		"https://medium.com/feed/tag/xss-attack",
		"https://medium.com/feed/tag/xss-vulnerability",
		"https://medium.com/feed/tag/xss-bypass",
		"https://medium.com/feed/tag/xss-filter-bypass",
		"https://medium.com/feed/tag/cross-site-scripting",
		"https://medium.com/feed/tag/dom-xss",
		"https://medium.com/feed/tag/blind-xss",
		"https://medium.com/feed/tag/stored-xss",
		"https://medium.com/feed/tag/owasp-top-10",
		"https://medium.com/feed/tag/owasp",
		"https://medium.com/feed/tag/sql-injection",
		"https://medium.com/feed/tag/web-application-security",
		"https://medium.com/feed/tag/injection",
		"https://medium.com/feed/tag/bug-bounty-writeup",
		"https://medium.com/feed/tag/vapt",
		"https://medium.com/feed/tag/vulnerability-assessment",
		"https://medium.com/feed/tag/cybersecurity",
		"https://medium.com/feed/tag/application-security",
		"https://medium.com/feed/tag/hacking",
		"https://medium.com/feed/tag/infosec",
		"https://medium.com/feed/tag/ctf",
		"https://medium.com/feed/tag/penetration-testing",
		"https://medium.com/feed/tag/writeup",
		"https://medium.com/feed/tag/tryhackme",
		"https://medium.com/feed/tag/vulnhub",
		"https://medium.com/feed/tag/security",
		"https://medium.com/feed/tag/bug-bounty",
		"https://medium.com/feed/tag/bug-hunter",
		"https://medium.com/feed/tag/golang",
		"https://medium.com/feed/tag/info-sec-writeup",
		"https://medium.com/feed/tag/hackthebox-writeup",
		"https://medium.com/feed/tag/ethical-hacking",
		"https://medium.com/feed/tag/api-security",
		"https://medium.com/feed/tag/hackerone",
		"https://medium.com/feed/tag/authentication",
		"https://medium.com/feed/tag/vulnerability",
		"https://medium.com/feed/tag/recon",
		"https://surya-dev.medium.com/feed",
		"https://infosecwriteups.com/feed",
		"https://medium.com/feed/@securitylit",
		"https://medium.com/feed/@tomnomnom",
		"https://medium.com/feed/@cappriciosec",
		"https://medium.com/feed/@302Found",
		"https://medium.com/feed/@newp_th",
		"https://medium.com/feed/@pdelteil",
		"https://ruvlol.medium.com/feed",
		"https://medium.com/@know.0nix/feed",
		"https://medium.com/@bugh4nter/feed",
		"https://seqrity.medium.com/feed",
		"https://vickieli.medium.com/feed",
		"https://medium.com/feed/intigriti",
		"https://medium.com/@intideceukelaire/feed",
		"https://medium.com/@projectdiscovery/feed",
		"https://jonathandata1.medium.com/feed",
		"https://medium.com/@Hacker0x01/feed",
		"https://medium.com/feed/pentesternepal",
		"https://0xjin.medium.com/feed",
		"https://medium.com/@infosecwriteups/feed",
		"https://medium.com/@jhaddix/feed",
		"https://medium.com/@NahamSec/feed",
		"https://orwaatyat.medium.com/feed",
		"https://zseano.medium.com/feed",
		"https://d0nut.medium.com/feed",
		"https://medium.com/feed/towards-aws",
		"https://medium.com/@stackzero/feed",
	}

	// Load previous data from file
	fileContents, err := ioutil.ReadFile("titles_guids.txt")
	if err != nil {
		fileContents = []byte{}
	}
	fileString := string(fileContents)

	// Discord webhook URL (replace with the actual URL)
	webhookURL := "Replace-With-Discord-webhook-URL"

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error fetching %s: %v", url, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to fetch %s (HTTP status: %d)", url, resp.StatusCode)
			continue
		}

		// Parse RSS feed
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("Error parsing feed %s: %v", url, err)
			continue
		}

		fmt.Println("Checking feed:", url)

		// Extract feed data
		doc.Find("item").Each(func(i int, s *goquery.Selection) {
			title := cleanCDATA(s.Find("title").Text())
			guid := cleanCDATA(s.Find("guid").Text())
			pubDate := cleanCDATA(s.Find("pubDate").Text())

			// Check if this article is new
			if !strings.Contains(fileString, guid) {
				// Append new data
				fileString += fmt.Sprintf("Title: %s\nGUID: %s\npubDate: %s\n\n", title, guid, pubDate)

				// Format message for Discord
				message := fmt.Sprintf("\n\n[+] Hey Nigga, There's a New writeup :peach:\n:pencil: [%s](%s)\n:clock1: %s", title, guid, pubDate)

				// Send to Discord
				sendToDiscord(webhookURL, message)
			}
		})
	}

	// Save updated data to file
	err = ioutil.WriteFile("titles_guids.txt", []byte(fileString), 0644)
	if err != nil {
		log.Printf("Error saving file: %v", err)
	}
}

// Clean <![CDATA[]]> from values
func cleanCDATA(text string) string {
	re := regexp.MustCompile("<!\\[CDATA\\[(.*)\\]\\]>")
	return re.ReplaceAllString(text, "$1")
}

func sendToDiscord(webhookURL, message string) {
	data := map[string]string{"content": message}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		return
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error creating webhook request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending message to Discord: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Discord webhook failed (HTTP status: %d): %s", resp.StatusCode, string(body))
	}
}

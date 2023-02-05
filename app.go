package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"bytes"
	"io/ioutil"
	"regexp"
	"github.com/PuerkitoBio/goquery"
	
)

func main() {
	// Define a list of URLs to scrape
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

	// Load the contents of the file, if it exists
	fileContents, err := ioutil.ReadFile("titles_guids.txt")
	if err != nil {
		fileContents = []byte{}
	}

	// Convert the contents of the file to a string
	fileString := string(fileContents)

	// Discord webhook URL
	webhookURL := "your-web-hook-URL"

	for _, url := range urls {
		// Load the URL
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("URL:", url)

		// Find the title and guid elements
		doc.Find("item").Each(func(i int, s *goquery.Selection) {
			title := s.Find("title").Text()
			title = regexp.MustCompile(`<!\[CDATA\[(.*)\]\]>`).ReplaceAllString(title, "$1")
			guid := s.Find("guid").Text()
			pubDate := s.Find("pubDate").Text()

			// Check if the title and GUID are in the file
			if !strings.Contains(fileString, title) && !strings.Contains(fileString, guid) && !strings.Contains(fileString, pubDate) {
				// If not, add them to the file
				fileString += fmt.Sprintf("Title: %s\nGUID: %s\n pubDate: %s\n" , title, guid, pubDate)

				// make hyperlink
				title = fmt.Sprintf("[%s](%s)", title, guid)

				// Send the new title and GUID to Discord
				payload := fmt.Sprintf("{\"content\": \"Hey Nigga, There's a New writeup :peach:  \\n%s\\n%s\"}", title, pubDate)
				req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer([]byte(payload)))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				_, err = client.Do(req)
				if err != nil {
					log.Fatal(err)
				}
			}

			fmt.Printf("\tTitle: %s\n\tGUID: %s\n", title, guid, pubDate)
		})
	}

	// Write the contents of the file back to disk
	err = ioutil.WriteFile("titles_guids.txt", []byte(fileString), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

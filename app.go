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
		"https://medium.com/feed/tag/ctf",
	}

	// Load the contents of the file, if it exists
	fileContents, err := ioutil.ReadFile("titles_guids.txt")
	if err != nil {
		fileContents = []byte{}
	}

	// Convert the contents of the file to a string
	fileString := string(fileContents)

	// Discord webhook URL
	webhookURL := "your-webhook-URL"

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

// Example code from https://www.makeuseof.com/parse-and-generate-html-in-go/

package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	// Send an HTTP GET request to the example.com web page
	resp, err := http.Get("https://www.chrisdunaj.com")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Use the html package to parse the response body from the request
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Find and print all links on the web page
	var links []string
	var link func(*html.Node)
	link = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					// adds a new link entry when the attribute matches
					links = append(links, a.Val)
				}
			}
		}

		// traverses the HTML of the webpage from the first child node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			link(c)
		}
	}
	link(doc)

	// loops through the links slice
	for _, l := range links {
		fmt.Println("Link:", l)
	}
}

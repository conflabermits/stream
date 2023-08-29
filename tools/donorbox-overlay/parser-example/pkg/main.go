// Example code from https://www.makeuseof.com/parse-and-generate-html-in-go/

package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/net/html"
)

func main() {

	//targetUrl := "http://localhost:8080/" // For local testing
	targetUrl := "https://donorbox.org/support-black-girls-code/fundraiser/christopher-dunaj" // For live testing

	fmt.Println("Fetching URL:", targetUrl)
	resp, err := http.Get(targetUrl)
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
	//var links []string
	var totalRaised float64
	var paidCount string
	var raiseGoal float64
	var link func(*html.Node)
	link = func(n *html.Node) {

		/* if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					// adds a new link entry when the attribute matches
					links = append(links, a.Val)
				}
			}
		} */

		dollarMatch, _ := regexp.MatchString("^\\$\\d{1,}", n.Data)

		if dollarMatch { //&& n.Type == html.ElementNode {
			for i := range (n.Parent).Attr {
				if (n.Parent).Attr[i].Val == "total-raised" {
					// Formatting the string to remove the dollar sign (https://www.makeuseof.com/go-formatting-numbers-currencies/)
					totalRaised, err = strconv.ParseFloat(n.Data[1:], 64)
					if err != nil {
						fmt.Println("Error:", err)
					}
				}
				if (n.Parent).Attr[i].Val == "bold" {
					raiseGoal, err = strconv.ParseFloat(n.Data[1:], 64)
					if err != nil {
						fmt.Println("Error:", err)
					}
				}
			}
		}

		numMatch, _ := regexp.MatchString("^\\d{1,}", n.Data)

		if n.Data != "" && n.Type == html.TextNode && numMatch {
			for i := range (n.Parent).Attr {
				if (n.Parent).Attr[i].Val == "paid-count" {
					paidCount = n.Data
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
	/* for _, l := range links {
		fmt.Println("Link:", l)
	} */
	fmt.Println("Number of contributors:", paidCount)
	fmt.Println("Total raised: $", totalRaised)
	fmt.Println("Raise goal: $", raiseGoal)
}

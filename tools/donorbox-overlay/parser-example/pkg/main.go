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
	// Send an HTTP GET request to the example.com web page
	resp, err := http.Get("http://localhost:8080/")
	//resp, err := http.Get("https://donorbox.org/support-black-girls-code/fundraiser/christopher-dunaj")
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

	/* Attrs from bs4
	progress = str(soup.find(attrs={"id": "panel-1"}).contents[1])
	num_donators = soup.find(attrs={"id": "paid-count"}).text
		#'0'
	total_raised = soup.find(attrs={"id": "total-raised"}).text
		#'$0'
	raise_goal = soup.find(attrs={"id": "panel-1"}).contents[1].find_all(attrs={"class": "bold"})[2].text
		#'$500'
	*/

	// Find and print all links on the web page
	var links []string
	//var panels []string
	//var dollarValues []string
	var totalRaised float64
	var paidCount string
	var raiseGoal float64
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

		dollarMatch, _ := regexp.MatchString("^\\$\\d{1,}", n.Data)

		if dollarMatch { //n.Type == html.ElementNode && regexpMatch {
			for i := range (n.Parent).Attr {
				if (n.Parent).Attr[i].Val == "total-raised" {
					//dollarValues = append(dollarValues, n.Data)
					//fmt.Println("Total raised:", n.Data)
					// Formatting the string to remove the dollar sign (https://www.makeuseof.com/go-formatting-numbers-currencies/)
					totalRaised, err = strconv.ParseFloat(n.Data[1:], 64)
					if err != nil {
						fmt.Println("Error:", err)
					}
					//fmt.Println("Total raised:", totalRaised)
				}
				if (n.Parent).Attr[i].Val == "bold" {
					raiseGoal, err = strconv.ParseFloat(n.Data[1:], 64)
					if err != nil {
						fmt.Println("Error:", err)
					}
					//fmt.Println("Total raised:", totalRaised)
				}
			}
			/* for _, div := range n.Attr {
				if div.Key == "description" {
					// adds a new link entry when the attribute matches
					panels = append(panels, div.Val)
				}
			} */
		}

		numMatch, _ := regexp.MatchString("^\\d{1,}", n.Data)
		if n.Data != "" && n.Type == html.TextNode && numMatch {
			for i := range (n.Parent).Attr {
				if (n.Parent).Attr[i].Val == "paid-count" {
					//dollarValues = append(dollarValues, n.Data)
					//fmt.Println("Total raised:", n.Data)
					// Formatting the string to remove the dollar sign (https://www.makeuseof.com/go-formatting-numbers-currencies/)
					paidCount = n.Data
					//fmt.Println("Total raised:", totalRaised)
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
	// loops through the panels slice
	/* for _, p := range panels {
		fmt.Println("Panel:", p)
	} */
	// loops through the dollarValues slice
	/* for _, d := range dollarValues {
		fmt.Println("Dollar values:", d)
	} */
	fmt.Println("Number of contributors:", paidCount)
	fmt.Println("Total raised: $", totalRaised)
	fmt.Println("Raise goal: $", raiseGoal)
}

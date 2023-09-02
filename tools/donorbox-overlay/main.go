package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"golang.org/x/net/html"
)

/*
TO DO:
* Turn HTML content into a template
* Give the HTML template some inputs for changing reload interval, maybe other things
* Give the HTML template a button to save changes and reload with new values
*/

type Options struct {
	Url  string
	Port string
}

var targetUrl string

func parseArgs() (*Options, error) {
	options := &Options{}
	flag.StringVar(&options.Url, "url", "http://localhost:8080", "Donorbox URL to check")
	flag.StringVar(&options.Port, "port", "38080", "Port to run the local web server")
	flag.Usage = func() {
		fmt.Printf("Usage: <app> [options]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	targetUrl = options.Url
	return options, nil
}

func main() {
	options, err := parseArgs()
	if err != nil {
		os.Exit(1)
	}
	http.HandleFunc("/", serveHTML)
	fmt.Printf("Server starting on http://localhost:" + options.Port + "\n")
	fmt.Printf("Server checking URL: " + options.Url + "\n")
	http.ListenAndServe(":"+options.Port, nil)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Auto-Reloading Web Page</title>
			<style type="text/css">
				.body {
					width: auto;
				}
				div.main {
					font-size: 48px;
					color: yellowgreen;
				}
				h1 {
					font-size: 64px;
					color: yellowgreen;
				}
			</style>
			<script>
				function reloadPage() {
					location.reload();
				}
				setTimeout(reloadPage, 61000); // Reload every N milliseconds
			</script>
		</head>
		<body>
			<h1>Donorbox progress:</h1>
			<div class="main">
				<p>` + getDonorboxProgress() + `</p>
			</div>
		</body>
		</html>
	`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlContent)
}

func getDonorboxProgress() string {
	//targetUrl := "http://localhost:8080/" // For local testing
	//targetUrl := "https://donorbox.org/support-black-girls-code/fundraiser/christopher-dunaj" // For live testing

	fmt.Println("Fetching URL:", targetUrl)
	resp, err := http.Get(targetUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return "Error"
	}

	defer resp.Body.Close()

	// Use the html package to parse the response body from the request
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "Error"
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
	fmt.Println("  Number of contributors:", paidCount)
	fmt.Printf("  Total raised: $%g\n", totalRaised)
	fmt.Printf("  Raise goal: $%g\n", raiseGoal)

	return fmt.Sprintf("Number of contributors: %s<BR>Total raised: $%g<BR>Raise goal: $%g", paidCount, totalRaised, raiseGoal)

}

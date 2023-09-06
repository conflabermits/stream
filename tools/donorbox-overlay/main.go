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
	Url     string
	Port    string
	Timeout int
}

var targetUrl string
var pageTimeout string

func parseArgs() (*Options, error) {
	options := &Options{}
	flag.StringVar(&options.Url, "url", "http://localhost:8080", "Donorbox URL to check")
	flag.StringVar(&options.Port, "port", "38080", "Port to run the local web server")
	flag.IntVar(&options.Timeout, "timeout", 60, "Page refresh rate, in seconds")
	flag.Usage = func() {
		fmt.Printf("Usage: <app> [options]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	targetUrl = options.Url
	pageTimeout = strconv.Itoa(options.Timeout * 1000)
	return options, nil
}

func main() {
	options, err := parseArgs()
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("Server starting on http://localhost:" + options.Port + "\n")
	fmt.Printf("Server checking URL: " + options.Url + "\n")

	http.HandleFunc("/", serveHTML)
	http.ListenAndServe(":"+options.Port, nil)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<link rel="icon" href="data:,">
			<title>Auto-Reloading Web Page</title>
			<style type="text/css">
				* {
					width: auto;
					font-family: Verdana, Arial, sans-serif;
					color: white;
					text-shadow: 0 0 2px blue, 0 0 4px hotpink;
				}
				h1 {
					font-size: 24px;
				}
				div.main {
					font-size: 18px;
				}
			</style>
			<script>
				function reloadPage() {
					location.reload();
				}
				setTimeout(reloadPage, ` + pageTimeout + `); // Reload every N milliseconds
			</script>
		</head>
		<body>
			<h1>Donorbox progress:</h1>
			<div class="main">
				<p><b>` + getDonorboxProgress() + `</b></p>
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

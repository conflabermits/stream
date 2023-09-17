# Donorbox Overlay

This is a simple overlay for Donorbox. It is designed to be used with broadcast software like Twitch Studio or OBS.

## Progress

This started as a simple Python script that scrapes the Donorbox page once and outputs the current progress.

```text
$ ./donorbox-scraper.py 
Number of donators: 0
Total raised: $0
Funraiser goal: $500
```

Later I was able to piece together a simple Go single page web server, `sample-html/sample.go`, that allows me to test out the web source overlay as a proof of concept, just to make sure I can clearly display content that updates regularly.

The latest development is an HTML scraper/parser written in Go that essentially does the same thing as a Python script: scrape the Donorbox page and output the current progress. Argument parsing has also been added to make it easier to specify the target URL to check and the local port to run the web server.

```text
$ go run main.go -port 38080 -timeout 15 -url https://donorbox.org/support-black-girls-code/fundraiser/christopher-dunaj
Server starting on http://localhost:38080
Server checking URL: https://donorbox.org/support-black-girls-code/fundraiser/christopher-dunaj
Fetching URL: https://donorbox.org/support-black-girls-code/fundraiser/christopher-dunaj
  Number of contributors: 2
  Total raised: $78.03
  Raise goal: $500
```

It's been integrated into an HTML single page web server so the page reload now calls the function and returns HTML-formatted text of the current fundraiser data.

![donorbox-overlay-html](donorbox-overlay-html.png "donorbox-overlay-html")

(The white background is treated as transparent in the web source overlay, and the white text shows up better against the background of brown wood paneling.)

## Next Steps

* Add text animations, sound effects, or other attention-grabbing features when the numbers change.
  * Define more CSS tags for alerts vs normal content.
  * Pass function content using CSS tags instead of hard-coded HTML tags.
  * Add logic in function to determine if old numbers match new numbers, and trigger additional alert content if not.
* Add logging.
* Write some tests.

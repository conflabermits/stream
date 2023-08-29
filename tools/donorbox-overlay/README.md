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

Later I was able to piece together a simple Go single page web server that allows me to test out the web source overlay as a proof of concept, just to make sure I can clearly display content that updates regularly. (That's what `main.go` currently does.)

The latest development is an HTML scraper/parser written in Go that essentially does the same thing as a Python script: scrape the Donorbox page and output the current progress.

```text
$ go run main.go 
Fetching URL: https://donorbox.org/support-black-girls-code/fundraiser/christopher-dunaj
Number of contributors: 1
Total raised: $ 53.03
Raise goal: $ 500
```

## Next Steps

Now that the HTML scraper/parser is written in Go, it should be easier to call it from the HTML web server. Tying these two things together would be the next functional step forward. Once they're successfully working together I can focus on cleaning up the HTML, adding logging and arg parsing, and writing some tests.

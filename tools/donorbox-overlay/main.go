package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	// Start a goroutine to reload the contents every 30 seconds
	contentReloadTicker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go reloadContents(contentReloadTicker, quit)

	// Start the web server
	http.HandleFunc("/", serveHTML)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func reloadContents(ticker *time.Ticker, quit chan struct{}) {
	for {
		select {
		case <-ticker.C:
			fmt.Println("Reloading contents...")
			// Here you can update the contents or fetch new data from external sources.
			// For this example, let's just print a message.
			fmt.Println("Contents reloaded!")
			// Here's some additional content, trying to be unique.
			fmt.Println("Random integer: ", rand.Int())
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	// Replace this with your actual HTML content if needed.
	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Auto-Reloading Web Page</title>
			<script>
				function reloadPage() {
					location.reload();
				}
				setTimeout(reloadPage, 5000); // Reload every N milliseconds
			</script>
		</head>
		<body>
			<h1>Hello, this is an auto-reloading web page!</h1>
            <h2>Here's some additional content, trying to be unique.</h2>
            <p>Random integer: ` + fmt.Sprint(rand.Int()) + `</p>
		</body>
		</html>
	`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlContent)
}

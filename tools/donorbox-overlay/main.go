package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

/*
TO DO:
* Add function for checking current donorbox values
* Turn HTML content into a template
* Give the HTML template some inputs for changing reload interval, maybe other things
* Give the HTML template a button to save changes and reload with new values
*/

func main() {
	http.HandleFunc("/", serveHTML)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
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
					font-size: 64px;
					color: yellowgreen;
				}
			</style>
			<script>
				function reloadPage() {
					location.reload();
				}
				setTimeout(reloadPage, 5000); // Reload every N milliseconds
			</script>
		</head>
		<body>
			<h1>Hello, this is an auto-reloading web page!</h1>
			<div class="main">
				<p>Random integer: ` + fmt.Sprint(rand.Int()) + `</p>
			</div>
		</body>
		</html>
	`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlContent)
}

## Summary

Quick note on how to create a poll.

## Command

```text
$ curl -X POST 'https://api.twitch.tv/helix/polls' -H 'Authorization: Bearer sdfyo32rsdf89uewfrf23ejbjkwdys' -H 'Client-Id: d9023rnjklsdncvsdo9iwojrnnsdff' -H 'Content-Type: application/json' -d '{
  "broadcaster_id":"123456789",
  "title":"Heads or Tails, but cats?",
  "choices":[{
    "title":"lion"
  },
  {
    "title":"tiger"
  }],
  "channel_points_voting_enabled":true,
  "channel_points_per_vote":1000,
  "duration":60
}'
```

## Response

```text
{"data":[{"id":"12345678-abcd-6543-wxyz-098765432199","broadcaster_id":"123456789","broadcaster_name":"conflabermits","broadcaster_login":"conflabermits","title":"Heads or Tails, but cats?","choices":[{"id":"98765432-xxxx-1111-zzzz-123456667899","title":"lion","votes":0,"channel_points_votes":0},{"id":"12121212-3434-5656-7878-9a8b7c6d5e4f","title":"tiger","votes":0,"channel_points_votes":0}],"channel_points_voting_enabled":true,"channel_points_per_vote":1000,"status":"ACTIVE","duration":60,"started_at":"2024-06-27T03:02:13.712731778Z"}]}
```

## Suggested POST code

Gemini suggested this as a simple straight-forward way to send a POST request in Go.

> ```text
> package main
> 
> import (
>   "bytes"
>   "encoding/json"
>   "fmt"
>   "io/ioutil"
>   "net/http"
> )
> 
> type MyData struct {
>   // Define your data structure here
>   Name string `json:"name"`
>   Value int `json:"value"`
> }
> 
> func main() {
>   // Define your URL and data
>   url := "https://your-api.com/endpoint"
>   data := MyData{Name: "Example Name", Value: 10}
> 
>   // Marshal the data to JSON
>   jsonData, err := json.Marshal(data)
>   if err != nil {
>     fmt.Println("Error marshalling JSON:", err)
>     return
>   }
> 
>   // Create a new POST request
>   req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
>   if err != nil {
>     fmt.Println("Error creating request:", err)
>     return
>   }
> 
>   // Set the Content-Type header
>   req.Header.Set("Content-Type", "application/json")
> 
>   // Send the request with a client
>   client := &http.Client{}
>   resp, err := client.Do(req)
>   if err != nil {
>     fmt.Println("Error sending request:", err)
>     return
>   }
>   defer resp.Body.Close()
> 
>   // Optionally, handle the response
>   body, err := ioutil.ReadAll(resp.Body)
>   if err != nil {
>     fmt.Println("Error reading response body:", err)
>     return
>   }
> 
>   fmt.Println("Response status:", resp.StatusCode)
>   fmt.Println("Response body:", string(body))
> }
> ```
> 
> This code does the following:
> 
> 1. Defines a struct MyData to represent your JSON data.
> 1. Creates a data object of type MyData with your actual values.
> 1. Marshals the data object into a byte array using json.Marshal.
> 1. Creates a new HTTP request object using http.NewRequest with the method set to "POST", the URL, and the JSON data as the request body.
> 1. Sets the Content-Type header to "application/json" to indicate the request body format.
> 1. Creates an http.Client and sends the request using client.Do.
> 1. Optionally, reads the response body and prints the status code and body content.

## Browser Source URL

https://www.twitch.tv/popout/conflabermits/poll

## References

* https://dev.twitch.tv/docs/api/reference/#create-poll


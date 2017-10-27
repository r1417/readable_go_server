package main

import (
  "fmt"
  "html"
  "log"
  "net/http"
)

func main() {
  log.Print("aaa");
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      log.Print("bbb");
      fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
  })

  log.Fatal(http.ListenAndServe(":8080", nil))
}

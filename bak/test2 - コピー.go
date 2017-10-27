package main

import (
  "fmt"
  "html"
  "log"
  "net/http"
  "time"
  "encoding/json"
  
  "github.com/julienschmidt/httprouter"

  )








func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
  fmt.Fprintf(w, "%q", html.EscapeString(ps.ByName("id")))
}


type Todo struct {
  Name      string    `json:"name"`
  Completed bool      `json:"completed"`
  Due       time.Time `json:"due"`
}

type Todos []Todo

func JTest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    todos := Todos{
        Todo{Name: "testName", Completed: false},
        Todo{Name: "testName2", Completed: true},
    }
    
    json.NewEncoder(w).Encode(todos)
}

func main() {
  router := httprouter.New()
  router.GET("/jtest", JTest)
  router.GET("/identify/:id", Index)

  log.Fatal(http.ListenAndServe(":8080", router))
}

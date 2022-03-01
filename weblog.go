package main

import (
//  "fmt"
  "net/http"
  "html/template"
)

type BlogPage struct {
  Title string
  Blog string
}

/*
func indexHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<h1>suck it bob</h1>")
}
*/

func blogHandler(w http.ResponseWriter, r *http.Request) {
  p := BlogPage{Title: "bob", Blog: "sucks ass"}
  t, _ := template.ParseFiles("template.html")
  t.Execute(w, p)
}

func main() {
  http.HandleFunc("/", blogHandler)
  http.ListenAndServe(":8000", nil)
}

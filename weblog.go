package main
// once this runs, it should already have the page consructed in data
// only one execution should be necessary

import (
  "net/http"
  "html/template"
  "os"
  "log"
//  "fmt"
)

/*
type BlogPage struct {
  Title string
  Blog string
}
*/

type ArchiveData struct {
  YearsBlogs [][]string
}

/*
func indexHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<h1>suck it bob</h1>")
}
*/

func blogHandler(w http.ResponseWriter, r *http.Request) {
  years, err := os.ReadDir("./blog")
  if err != nil {
      log.Fatal(err)
  }
  yearsBlogs := make([][]string, len(years))
  i:=0
  for _, y := range years {
    files, err := os.ReadDir("./blog/" + y.Name())
    if err != nil {
        log.Fatal(err)
    }
    yearsBlogs[i] = make([]string, len(files))
    yearsBlogs[i] = append(yearsBlogs[i], y.Name())
    for _, f := range files {
      yearsBlogs[i] = append(yearsBlogs[i], f.Name())
    }
    i = i + 1
  }
  p := ArchiveData{YearsBlogs: yearsBlogs}
  t, _ := template.ParseFiles("blog.html")
  t.Execute(w, p)
/*
  filename := "carl.html"
  body, _ := template.ParseFiles(filename)
  body.Execute(w, p)
*/
}

func main() {
  http.HandleFunc("/", blogHandler)
  http.ListenAndServe(":8000", nil)
}


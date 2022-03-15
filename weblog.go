package main
// once this runs, it should already have the page consructed in data
// only one execution should be necessary

import (
  "net/http"
  "html/template"
  "os"
  "log"
  "fmt"
)

/*
type BlogPage struct {
  Title string
  Blog string
}
*/

type BlogHome struct {
  Years []string
  Blogs []string
}

/*
func indexHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<h1>suck it bob</h1>")
}
*/

func blogHandler(w http.ResponseWriter, r *http.Request) {
  var blogYears = []string{}
  var blogs = []string{}
  files, err := os.ReadDir("./blog")
  if err != nil {
      log.Fatal(err)
  }
  for _, f := range files {
    blogYears = append(blogYears, f.Name())
    giles, err := os.ReadDir("./blog/" + f.Name())
    if err != nil {
        log.Fatal(err)
    }
    for _, g := range giles {
      blogs = append(blogs, g.Name())
      fmt.Println(blogs)
      fmt.Fprintf(w, blogs)
    }
    blogs = nil
  }

  p := BlogHome{Years: blogYears, Blogs: blogs}
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


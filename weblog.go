/*package main

// once this runs, it should already have the page consructed in data
// only one execution should be necessary

import (
	"html/template"
	"log"
	"net/http"
	"os"
	//  "fmt"
)

/*
type BlogPage struct {
  Title string
  Blog string
}
*/

const templateDir string = "./template/"
const staticDir string = "./static/"
const blogDir string = "./blog/"

type ArchiveData struct {
	YearsBlogs [][]string
}

/*
func indexHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<h1>suck it bob</h1>")
}
*/

func bang(w http.ResponseWriter, r *http.Request) {
	yearsBlogs := make([][]string, 5)
	p := ArchiveData{YearsBlogs: yearsBlogs}
	t, _ := template.ParseFiles(blogDir + "2021/12-24_christmas-eve.html")
	t.Execute(w, p)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	years, err := os.ReadDir(blogDir)
	if err != nil {
		log.Fatal(err)
	}
	yearsBlogs := make([][]string, len(years))
	i := 0
	for _, y := range years {
		files, err := os.ReadDir(blogDir + y.Name())
		if err != nil {
			log.Fatal(err)
		}
		yearsBlogs[i] = []string{}
		yearsBlogs[i] = append(yearsBlogs[i], y.Name())
		for _, f := range files {
			yearsBlogs[i] = append(yearsBlogs[i], f.Name())
		}
		i = i + 1
	}
	p := ArchiveData{YearsBlogs: yearsBlogs}
	t, _ := template.ParseFiles(templateDir + "blog.html")
	t.Execute(w, p)
	/*
	   filename := "carl.html"
	   body, _ := template.ParseFiles(filename)
	   body.Execute(w, p)
	*/
}

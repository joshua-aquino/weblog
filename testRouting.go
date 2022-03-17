package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func addRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

var routes = []route{
	addRoute("GET", "/", home),
	addRoute("GET", "/blog", blog),
}

//------------------------------------------------------------------------------

const templateDir string = "./template/"
const staticDir string = "./static/"
const blogDir string = "./blog/"

type Archive struct {
	YearsBlogs [][]string
}

func home(w http.ResponseWriter, r *http.Request) {
	p := ctxKey{}
	t, _ := template.ParseFiles(templateDir + "home.html")
	t.Execute(w, p)
}

func blog(w http.ResponseWriter, r *http.Request) {
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
	p := Archive{YearsBlogs: yearsBlogs}
	t, _ := template.ParseFiles(templateDir + "blog.html")
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", Serve)
	http.ListenAndServe(":8000", nil)
}

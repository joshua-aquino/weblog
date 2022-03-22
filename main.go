package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var LayoutDir string = "views/layouts/"

func NewView(layout string, files ...string) *View {
	files = append(layoutFiles(), files...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*.tmpl")
	if err != nil {
		panic(err)
	}
	return files
}

//-----------------------------------------------------------------------------

type ArchiveEntry struct {
	Title, FileName, BlogDir string
	IsYear                   bool
}

// the title finding function should not be run everytime you update the page,
// it should be run on a loop in a different thread every hour or so
func completeArchives() []ArchiveEntry {
	var archiveEntries []ArchiveEntry
	archiveEntry := ArchiveEntry{}
	lastYear := "XX"
	re := regexp.MustCompile(`<title>[^<]*`)
	files, err := os.ReadDir(blogDir)
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	for _, f := range files {
		archiveEntry = ArchiveEntry{f.Name(), f.Name(), blogDir, false}
		presentYear := strings.Split(f.Name(), "-")[0]
		if presentYear != lastYear {
			archiveEntry.Title = presentYear
			archiveEntry.IsYear = true
			lastYear = presentYear
			archiveEntries = append(archiveEntries, archiveEntry)
			i = i + 1
		}
		archiveEntry.IsYear = false
		readFile, err := os.Open(blogDir + f.Name())
		if err != nil {
			log.Fatalf("failed to open file: %s", err)
		}
		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)
		var fileTextLines []string
		for fileScanner.Scan() {
			fileTextLines = append(fileTextLines, fileScanner.Text())
		}
		readFile.Close()
		for _, eachline := range fileTextLines {
			if re.MatchString(eachline) {
				archiveEntry.Title =
					strings.ReplaceAll(re.FindString(eachline), "<title>", "")
				continue
			}
		}
		archiveEntries = append(archiveEntries, archiveEntry)
		i = i + 1
	}
	return archiveEntries
}

//-----------------------------------------------------------------------------

const viewsDir string = "views/"
const staticDir string = "static/"
const blogDir string = "blogs/"

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
	addRoute("GET", "/", indexHandler),
	addRoute("GET", "/blogs(/)?", archiveHandler),
	addRoute("GET", "/blogs/([^/]+)", blogHandler),
}

//-----------------------------------------------------------------------------

var index *View
var archive *View
var blog *View

func main() {
	index = NewView("default", "views/index.gohtml")
	archive = NewView("default", "views/archive.gohtml")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", Serve)
	http.ListenAndServe(":8000", nil)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	index.Render(w, nil)
}

func archiveHandler(w http.ResponseWriter, r *http.Request) {
	p := completeArchives()
	archive.Render(w, p)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	blog = NewView("default", "blogs/"+getField(r, 0))
	blog.Render(w, nil)
}

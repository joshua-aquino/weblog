package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const templateDir string = "./template/"
const staticDir string = "./static/"
const blogDir string = "./blog/"

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
	addRoute("GET", "/blog(/)?", blogArchive),
	addRoute("GET", "/blog/([^/]+)", blog),
}

//------------------------------------------------------------------------------

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
		archiveEntry = ArchiveEntry{"none", "none", blogDir, false}
		presentYear := strings.Split(f.Name(), "-")[0]

		if presentYear != lastYear {
			//make a new entry and don't do the rest
			archiveEntry.Title = presentYear
			archiveEntry.IsYear = true
			lastYear = presentYear
			archiveEntries = append(archiveEntries, archiveEntry)
			i = i + 1
		}
		archiveEntry.IsYear = false
		var titleFound bool = false
		archiveEntry.FileName = f.Name()
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
				titleFound = true
			}
		}
		if !titleFound {
			archiveEntry.Title = archiveEntry.FileName
		}
		archiveEntries = append(archiveEntries, archiveEntry)
		i = i + 1
	}
	return archiveEntries
}

//------------------------------------------------------------------------------

func home(w http.ResponseWriter, r *http.Request) {
	p := ctxKey{}
	t, _ := template.ParseFiles(templateDir + "home.html")
	t.Execute(w, p)
}

func blogArchive(w http.ResponseWriter, r *http.Request) {
	p := completeArchives()
	log.Println(p[0].Title)
	t, _ := template.ParseFiles(templateDir + "archive.html")
	t.Execute(w, p)
	// template files cannot read the individual values within struct arrays; I
	// will try to make an in-file template next
}

func blog(w http.ResponseWriter, r *http.Request) {
	p := ctxKey{}
	t, _ := template.ParseFiles(blogDir + getField(r, 0))
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", Serve)
	http.ListenAndServe(":8000", nil)
}

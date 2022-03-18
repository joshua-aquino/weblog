package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

const templateDir string = "./template/"
const staticDir string = "./static/"
const blogDir string = "./blog/"

var titles = []string{}

type Data struct{}

func findBlogTitles(inputFiles []string) {
	re := regexp.MustCompile(`<title>[^<]*`)
	files, err := os.ReadDir(blogDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		var titleFound bool = false
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
				titles = append(titles,
					strings.ReplaceAll(re.FindString(eachline), "<title>", ""))
				titleFound = true
			}
		}
		if !titleFound {
			titles = append(titles, f.Name())
		}
	}
	log.Println("---------------------------------------")
	i := 0
	lastYear := strings.Split(inputFiles[i], "-")[0]
	log.Println("20" + strings.Split(inputFiles[i], "-")[0])
	for _, title := range titles {
		if strings.Split(inputFiles[i], "-")[0] != lastYear {
			log.Println("20" + strings.Split(inputFiles[i], "-")[0])
			lastYear = strings.Split(inputFiles[i], "-")[0]
		}
		log.Println(title)
		log.Println(inputFiles[i])
		i = i + 1
	}
}

func main() {
	var fileParts = []string{}
	var fileNames = []string{}
	var fileDates = []string{}
	files, err := os.ReadDir(blogDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fileParts = strings.Split(f.Name(), "_")
		fileNames = append(fileNames, f.Name())
		fileDates = append(fileDates, fileParts[0])
	}

	fileCount := len(files)
	log.Println(fileCount)
	log.Println(fileNames)
	log.Println(fileDates)
	findBlogTitles(fileNames)
}

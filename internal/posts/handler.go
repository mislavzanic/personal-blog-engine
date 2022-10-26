package posts

import (
	"io"
	"fmt"
	"log"
	"time"
	"bufio"
	"bytes"
	"strings"
	"net/http"
	"io/ioutil"
	"html/template"

	"github.com/gorilla/mux"
)

type Page struct {
	Body     string
	MetaData MetaData
	Path     string
}

type Posts struct {
	Pages  []*Page
}

type MetaData struct {
	Title string
	Date  time.Time
	Tags  []string
}


func PageHandler(w http.ResponseWriter, req *http.Request) {
	pageId := mux.Vars(req)["pageId"]
	p := readBlogPost(fmt.Sprintf("posts/%s.md", pageId))

	renderFromTemplate(w, "post.html", "html/post.html", template.FuncMap{"markDown": markDowner}, p)
}


func ViewAllPosts(w http.ResponseWriter, req *http.Request) {
	files, err := ioutil.ReadDir("posts")

	if err != nil {
		log.Fatal(err)
	}

	posts := Posts{}

	for _, file := range files {
		p := readBlogPost(fmt.Sprintf("posts/%s", file.Name()))
		posts.Pages = append(posts.Pages, p)
	}

	renderFromTemplate(w, "index.html", "html/index.html", template.FuncMap{"toURL": getUrl}, posts)
}

func readBlogPost(path string) *Page {
	body, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	metaData := readMetadata(bytes.NewReader(body))

	return &Page{MetaData: metaData, Body: readBody(body), Path: path}
}

func readMetadata(r io.Reader) MetaData {
	metaData := make([]string, 0)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			break
		}
		metaData = append(metaData, line)
	}
	date, err := time.Parse("2006-01-02", metaData[1])

	if err != nil {
		log.Fatal(err)
	}

	return MetaData{Title: metaData[0], Date: date, Tags: strings.Split(metaData[2], ",")}
}

func readBody(byteArray []byte) string {
	return string(bytes.Split(byteArray, []byte("---\n"))[1])
}

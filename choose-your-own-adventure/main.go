package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var tmpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Choose your own adventure</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		<p>
			{{range .Story}}{{.}}{{end}}
		</p>
		
		{{if .Options}}
			<p><b>Choose the next chapter you prefer:</b></p>
			<ul style="list-style: none;">
			{{range .Options}}
				<li>
					{{.Text}} <a href="/{{.Chapter}}">&rarr;</a>
				</li>
			{{end}}
			</ul>
		{{else}}
			<h3>END</h3>
			<img src="https://golang.org/doc/gopher/ref.png" />
		{{end}}
	</body>
</html>`

type Story map[string]Chapter

// https://mholt.github.io/json-to-go/
type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func main() {
	bytes, err := ioutil.ReadFile("gopher.json")
	if err != nil {
		log.Println("Error reading gopher.json", err)
		return
	}

	chapter, err := parseJSON(bytes)
	if err != nil {
		log.Println("Error parsing JSON", err)
		return
	}

	t := template.New("page")
	t, err = t.Parse(tmpl)
	if err != nil {
		log.Println("Error parsing template", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		page := "intro"
		if path != "/" {
			page = path[1:]
		}
		t.Execute(w, chapter[page])
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseJSON(b []byte) (Story, error) {
	chapter := Story{}
	err := json.Unmarshal(b, &chapter)
	if err != nil {
		return chapter, err
	}

	return chapter, nil
}

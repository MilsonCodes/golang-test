package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html", "templates/all.html"))
var validPath = regexp.MustCompile("^/((edit|save|view)/([a-zA-Z0-9]+)|all/)$")

// File ... To be displayed as a page with a title and body
type File struct {
	Name string
	Body []byte
}

func getFolderContents(dirName string) []*File {
	dir, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	ret := make([]*File, len(dir))

	for i, file := range dir {
		if file.IsDir() {
			// getFolderContents(file.Name(), depth+1)
			ret[i] = &File{Name: "FOLDER", Body: []byte(file.Name())}
		} else {
			ret[i], err = loadPage(strings.Split(file.Name(), ".")[0])
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// fmt.Println(ret)

	return ret
}

// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
// 	m := validPath.FindStringSubmatch(r.URL.Path)
// 	if m == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("invalid Page Title")
// 	}
// 	return m[2], nil
// }

func (p *File) save() error {
	filename := "files/" + p.Name + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*File, error) {
	filename := "files/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &File{Name: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func viewAllHandler(w http.ResponseWriter, r *http.Request, title string) {
	dir := getFolderContents("./files")
	renderTemplate(w, "all", dir)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &File{Name: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &File{Name: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/all/", makeHandler(viewAllHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

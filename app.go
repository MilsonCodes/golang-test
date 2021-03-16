package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Base directory for file storage
const rootDir = "./files/"

// Some extensions
var ext = []string{".txt", ".md", ".html"}

// Parse templates from templates folder
var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html", "templates/all.html"))

// Restrict allowed path
var validPath = regexp.MustCompile("^/((edit|save|view)/([a-zA-Z0-9]+))|all/$")

// File ... To be displayed as a page with a title and body
type File struct {
	Name string
	Body []byte
}

// Return the contents of the folder called dirName
// TODO: Handle recursively exploring folders
func getFolderContents(dirName string) []*File {
	dir, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	ret := make([]*File, len(dir))

	for i, file := range dir {
		if file.IsDir() {
			// getFolderContents(file.Name(), depth+1)
			ret[i] = &File{Name: "FOLDER", Body: []byte(file.Name() + "/")} // Temporarily store name in body of file struct
		} else {
			fileName := fmt.Sprintf(strings.Split(file.Name(), ".")[0]) // Get file name, ignore extension
			// TODO: Handle multiple extensions
			ret[i], err = loadPage(fileName)
			if err != nil {
				ret[i].print()
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

// Save file, write to system
func (p *File) save() error {
	filename := rootDir + p.Name + ext[0]
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Print information to the console
func (p *File) print() error {
	_, err := fmt.Printf("Name: %s\nBody: %s\n", p.Name, p.Body)
	return err
}

// Read file from the system, then return a struct of the file
func loadPage(title string) (*File, error) {
	filename := rootDir + title + ext[0]
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &File{Name: title, Body: body}, nil
}

/*
* Render the template based on input string and html extension
* Use p to pass in any data you want to be used in the template
 */
func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+ext[2], p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler wrapper to handle errors and validate path (reduce code)
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		// fmt.Println(m)
		fn(w, r, m[3])
	}
}

func viewAllHandler(w http.ResponseWriter, r *http.Request, title string) {
	dir := getFolderContents(rootDir)
	renderTemplate(w, "all", dir)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	// fmt.Println(title)
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

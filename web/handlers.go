package web

import (
	"fmt"
	"github.com/SpencerBrown/mongodb-repro/staticContent"
	"html/template"
	"log"
	"net/http"
)

var verbose bool // yes it's a global variable, deal with it

func initHandlers(m *http.ServeMux) {
	m.HandleFunc("/favicon.ico", handleIcon)
	m.HandleFunc("/", handleInitial)
	m.HandleFunc("/changes", handleChanges)
}

// Handle request for the icon, just say "not found"
func handleIcon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	logRequest("favicon", r)
}

// Handle / by putting up a form with the sctruct
func handleInitial(w http.ResponseWriter, r *http.Request) {
	logRequest("Initial", r)
	var formData string
	runTemplate(staticContent.Struct_html, "requestForm", formData, w)
}

// Handle /changes when user submits the form
func handleChanges(w http.ResponseWriter, r *http.Request) {
	var err error

	if err = r.ParseForm(); err != nil {
		fmt.Fprintf(w, "error parsing form: %v", err)
		return
	}
	logRequest("Changes", r)

	_, ok := r.PostForm["quit"]
	if ok {
		http.Redirect(w, r, "/shutdown", http.StatusMovedPermanently)
		return
	}
}

func runTemplate(templateString string, templateName string, input interface{}, w http.ResponseWriter) {
	templ, err := template.New(templateName).Parse(templateString)
	if err != nil {
		log.Fatalf("Internal error creating template %s: %v\n", templateName, err)
	}
	if err = templ.Execute(w, input); err != nil {
		log.Fatalf("Internal error executing template %s: %v\n", templateName, err)
	}
}

func logRequest(text string, r *http.Request) {
	if verbose {
		url := r.URL.String()
		values := r.Form.Encode()
		log.Printf("%s URL: %s Values: %s\n", text, url, values)
	}
}

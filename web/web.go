package web

import (
	"context"
	"github.com/SpencerBrown/mongodb-repro/staticContent"
	"github.com/zserge/lorca"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var verbose bool // yes it's a global variable, deal with it

func Server(host string, port int, v bool) {
	hostport := host + ":" + strconv.Itoa(port)
	verbose = v
	var ui lorca.UI = nil
	var err error
	ui, err = lorca.New("http://"+hostport, "", 1024, 1024)
	if err != nil {
		log.Fatalf("lorca error: %v\n", err)
	}
	defer ui.Close()
	//<-ui.Done()
	log.Printf("Starting Web server on %s (use Ctrl+C to terminate)\n", hostport)
	m := http.NewServeMux()
	m.HandleFunc("/favicon.ico", handleIcon)
	m.HandleFunc("/", handleRequest)
	s := http.Server{Addr: hostport, Handler: m}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		_, err = w.Write([]byte("OK"))
		// Cancel the context on request
		cancel()
	})
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	if ui == nil {
		<-ctx.Done()
	} else {
		select {
		case <-ctx.Done():
			{ // Shutdown the server when the context is canceled
				err = s.Shutdown(ctx)
			}
		case <-ui.Done():
			{
				err = s.Shutdown(ctx)
			}
		}
	}
	log.Print("Finished")
}

// Handle request for the icon, just say "not found"
func handleIcon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	logRequest("favicon", r)
}

// Handle / by putting up a form
func handleRequest(w http.ResponseWriter, r *http.Request) {
	logRequest("Initial", r)
	var formData string
	runTemplate(staticContent.Struct_html, "requestForm", formData, w)
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

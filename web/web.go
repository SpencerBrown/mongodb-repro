package web

import (
	"context"
	"github.com/zserge/lorca"
	"log"
	"net/http"
	"strconv"
)

func Server(host string, port int, v bool) {
	hostport := host + ":" + strconv.Itoa(port)
	verbose = v
	var ui lorca.UI = nil
	var err error
	log.Printf("Starting Web server on %s (use Ctrl+C to terminate)\n", hostport)
	m := http.NewServeMux()
	s := http.Server{Addr: hostport, Handler: m}
	initHandlers(m)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		logRequest("Quit", r)
		_, err = w.Write([]byte("OK"))
		// Cancel the context on request
		cancel()
	})
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	ui, err = lorca.New("http://"+hostport, "", 1024, 1024)
	if err != nil {
		log.Fatalf("lorca error: %v\n", err)
	}
	defer ui.Close()
	//<-ui.Done()
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

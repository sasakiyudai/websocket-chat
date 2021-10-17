package main

import (
	"flag"
	"log"
	"net/http"
	// "os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/sasakiyudai/websocket-chat/room"
	// "github.com/sasakiyudai/websocket-chat/trace"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "address of app")
	flag.Parse()
	r := room.NewRoom()
	// r.Tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.Run()
	log.Println("web server started on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

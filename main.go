package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"

	"github.com/sasakiyudai/websocket-chat/auth"
	"github.com/sasakiyudai/websocket-chat/room"
	"github.com/sasakiyudai/websocket-chat/trace"
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
	var f = flag.Bool("f", false, "trace will be displaied")
	flag.Parse()
	// OAuth

	r := room.NewRoom()
	if *f == true {
		r.Tracer = trace.New(os.Stdout)
	}
	http.Handle("/chat", auth.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", auth.LoginHandler)
	http.Handle("/room", r)
	go r.Run()
	log.Println("web server started on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

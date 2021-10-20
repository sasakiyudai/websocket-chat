package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"

	"github.com/sasakiyudai/websocket-chat/auth"
	"github.com/sasakiyudai/websocket-chat/room"
	"github.com/sasakiyudai/websocket-chat/trace"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf(".env: %v", err)
	}	
}

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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "address of app")
	var f = flag.Bool("f", false, "trace will be displaied")
	flag.Parse()
	// OAuth
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		facebook.New(os.Getenv("FACEBOOK_API_CLIENT_ID"), os.Getenv("FACEBOOK_API_SECRET_KEY"), "http://localhost:8080/auth/callback/facebook"),
		google.New(os.Getenv("GOOGLE_API_CLIENT_ID"), os.Getenv("GOOGLE_API_SECRET_KEY"), "http://localhost:8080/auth/callback/google"),
		github.New(os.Getenv("GITHUB_API_CLIENT_ID"), os.Getenv("GITHUB_API_SECRET_KEY"), "http://localhost:8080/auth/callback/github"),
	)
	r := room.NewRoom()
	if *f == true {
		r.Tracer = trace.New(os.Stdout)
	}
	http.Handle("/chat", auth.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", auth.LoginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room", r)
	go r.Run()
	log.Println("web server started on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

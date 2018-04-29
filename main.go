package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	r := newRoom()
	gomniauth.SetSecurityKey("NASZ KLUCZ AUTORYZACYJNY")
	gomniauth.WithProviders(
		google.New("identyfikator", "klucz_tajny",
			"http://localhost:8080/auth/callback/google"),
	)
	http.Handle("/chat", AuthRoute(&templateHandler{filename: "chat.html"}))

	http.Handle("/room", r)

	http.Handle("/login", &templateHandler{filename: "login.html"})

	http.HandleFunc("/auth/", loginHandler)

	go r.run()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Fatal:", err)
	}
}

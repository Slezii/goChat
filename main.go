package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	r := newRoom()
	gomniauth.SetSecurityKey("QQQQQBBBBBAAAAA")
	gomniauth.WithProviders()
	http.Handle("/chat", AuthRoute(&templateHandler{filename: "chat.html"}))

	http.Handle("/room", r)

	http.Handle("/login", &templateHandler{filename: "login.html"})

	http.HandleFunc("/auth/", loginHandler)

	go r.run()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Fatal:", err)
	}
}

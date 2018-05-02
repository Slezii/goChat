package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	. "./configs"
	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
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

	conf := config.NewConfig()
	conf.Load(file.NewSource(
		file.WithPath("config.json"),
	))
	var oAuthConf OAuthConfig
	conf.Get("hosts", "googleOauth").Scan(&oAuthConf)

	var serverConf ServerConfig
	conf.Get("hosts", "server").Scan(&serverConf)
	gomniauth.SetSecurityKey("QQQQQBBBBBAAAAA")
	gomniauth.WithProviders(
		google.New(oAuthConf.Id, oAuthConf.Secret,
			"http://"+serverConf.Address+":"+serverConf.Port+"/auth/callback/google"),
	)
	http.Handle("/chat", AuthRoute(&templateHandler{filename: "chat.html"}))

	http.Handle("/room", r)

	http.Handle("/login", &templateHandler{filename: "login.html"})

	http.HandleFunc("/auth/", loginHandler)

	go r.run()
	if err := http.ListenAndServe(":"+serverConf.Port, nil); err != nil {
		log.Fatal("Fatal:", err)
	}
}

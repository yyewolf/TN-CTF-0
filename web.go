package main

import (
	"embed"
	"encoding/gob"
	"fmt"
	"net/http"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/gorilla/pat"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// This is used to access the FS 'www', and compile the files in the binary.
//go:embed www
var webbox embed.FS

// Cookie is just storing username to allow the sql injection.
type cookieUser struct {
	Username      string
	Authenticated bool
}

// We will verify emails possible existence but nothing else.
var (
	verifier = emailverifier.
		NewVerifier().
		EnableSMTPCheck().
		EnableAutoUpdateDisposable()
)

var store *sessions.CookieStore

func hostService() {
	mux := pat.New()
	srv := http.Server{
		Addr:    ":1337",
		Handler: mux,
	}
	// We make sure cookies are secure.
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   0,
		HttpOnly: true,
	}

	gob.Register(cookieUser{})

	// ROUTING
	mux.Post("/login", compressHandler(http.HandlerFunc(login)))
	mux.Post("/signup", compressHandler(http.HandlerFunc(signup)))
	mux.Post("/transfer", compressHandler(http.HandlerFunc(webTransfer)))

	mux.Get("/login", compressHandler(http.HandlerFunc(login)))
	mux.Get("/signup", compressHandler(http.HandlerFunc(signup)))
	mux.Get("/logout", compressHandler(http.HandlerFunc(logout)))
	mux.Get("/", compressHandler(http.HandlerFunc(indexHandler)))

	go srv.ListenAndServe()
}

func (u *User) formatConnected(data string) string {
	// Simple formatting functions.
	data = strings.ReplaceAll(data, "{Error}", "")
	data = strings.ReplaceAll(data, "{User.Email}", u.Email)
	data = strings.ReplaceAll(data, "{User.Username}", u.Username)

	data = strings.ReplaceAll(data, "{User.PartieEntiere}", fmt.Sprint(u.PartieEntiere))
	data = strings.ReplaceAll(data, "{User.PartieDecimale}", fmt.Sprint(u.PartieDecimale))
	return data
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	datb, _ := webbox.ReadFile("www/index.html")
	dat := string(datb)
	Path := r.URL.Path[1:]
	Path = "www/" + Path
	if Path == "www/" {
		// MAIN INDEX PAGE DISPLAY HERE
		fmt.Fprint(w, dat)
	} else {
		// We probably don't need support for JS, SVG and ICO but who would mind :p
		dat, err := webbox.ReadFile(Path)
		if err != nil {
			fmt.Fprint(w, err)
		}
		if strings.HasSuffix(Path, ".css") {
			w.Header().Add("Content-Type", "text/css")
		}
		if strings.HasSuffix(Path, ".js") {
			w.Header().Add("Content-Type", "text/javascript")
		}
		if strings.HasSuffix(Path, ".svg") {
			w.Header().Add("Content-Type", "image/svg+xml")
		}
		if strings.HasSuffix(Path, ".ico") {
			w.Header().Add("Content-Type", "image/x-icon")
		}
		fmt.Fprint(w, string(dat))
	}
}

func webTransfer(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Redirect(w, r, "https://ctf0.yewolf.ovh/login", http.StatusSeeOther)
	}
	var val interface{}
	var ok bool
	if val, ok = session.Values["user"]; !ok {
		http.Redirect(w, r, "https://ctf0.yewolf.ovh/login", http.StatusSeeOther)
		return
	}

	cookieU := val.(cookieUser)
	user, err := getUserByUsername(cookieU.Username)
	if err != nil {
		datb, _ := webbox.ReadFile("www/login.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox(err))
		fmt.Fprint(w, dat)
		return
	}

	sendTo := r.FormValue("username")
	amount := r.FormValue("amount")

	// Avoids cheats by just sharing to your friends the money but if this happens then they could just share the trick :/
	if user.Username != "admin" {
		datb, _ := webbox.ReadFile("www/connected.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox("Your account is too recent to be able to transfer."))
		dat = user.formatConnected(dat)
		fmt.Fprint(w, dat)
		return
	}

	if sendTo == "" {
		datb, _ := webbox.ReadFile("www/connected.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox("You cannot send to nobody."))
		dat = user.formatConnected(dat)
		fmt.Fprint(w, dat)
		return
	}

	if amount == "" {
		datb, _ := webbox.ReadFile("www/connected.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox("Amount can't be empty."))
		dat = user.formatConnected(dat)
		fmt.Fprint(w, dat)
		return
	}

	sendingE, sendingD, err := textAmountToVal(amount)
	if err != nil {
		datb, _ := webbox.ReadFile("www/connected.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox(err))
		dat = user.formatConnected(dat)
		fmt.Fprint(w, dat)
		return
	}

	// We finally do the transfer.
	err = user.transfer(sendTo, sendingE, sendingD)
	if err != nil {
		datb, _ := webbox.ReadFile("www/connected.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox(err))
		dat = user.formatConnected(dat)
		fmt.Fprint(w, dat)
		return
	}

	datb, _ := webbox.ReadFile("www/connected.html")
	dat := user.formatConnected(string(datb))
	dat = strings.ReplaceAll(dat, "{Error}", errorBox("Your money has been transferred."))
	fmt.Fprint(w, dat)
}

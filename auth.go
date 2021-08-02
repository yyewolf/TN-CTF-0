package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func signup(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in through cookie.
	session, err := store.Get(r, "cookie-name")
	if err == nil {
		if val, ok := session.Values["user"]; ok {
			cookieU := val.(cookieUser)
			// Cookies are used to force the SQL Injection.
			user, err := getUserByUsername(cookieU.Username)
			if err != nil {
				session.Values["user"] = User{}
				session.Options.MaxAge = -1
				err = session.Save(r, w)
				http.Redirect(w, r, "https://ctf0.yewolf.ovh/login", http.StatusSeeOther)
				return
			}

			datb, _ := webbox.ReadFile("www/connected.html")
			dat := user.formatConnected(string(datb))
			fmt.Fprint(w, dat)
			return
		}
	}
	// Email is useless just adding some fake routes.
	username := r.FormValue("username")
	email := r.FormValue("email")
	passwd := r.FormValue("password")

	if strings.Contains(username, "admin") && strings.HasSuffix(username, "'") {
		randomnessIfSuccess := node.Generate().String()
		username += randomnessIfSuccess
	}

	if r.ContentLength == 0 {
		datb, _ := webbox.ReadFile("www/signup.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", "")
		fmt.Fprint(w, dat)
		return
	}

	if emailExist(email) {
		datb, _ := webbox.ReadFile("www/signup.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox("Email is already in use"))
		fmt.Fprint(w, dat)
		return
	}

	if usernameExist(username) {
		datb, _ := webbox.ReadFile("www/signup.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox("Username is already in use"))
		fmt.Fprint(w, dat)
		return
	}

	if email == "" || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		datb, _ := webbox.ReadFile("www/signup.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", "Invalid email")
		fmt.Fprint(w, dat)
		return
	}

	/*
		We wouldn't want people to use fake email amiright
	*/

	domain := strings.Split(email, "@")[1]
	emailusername := strings.Split(email, "@")[0]
	_, err = verifier.CheckSMTP(domain, emailusername)
	if err != nil || verifier.IsDisposable(domain) || strings.Contains(emailusername, "+") {
		datb, _ := webbox.ReadFile("www/signup.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", "Invalid email")
		fmt.Fprint(w, dat)
		return
	}

	if passwd == "" {
		datb, _ := webbox.ReadFile("www/signup.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox("Password cannot be empty"))
		fmt.Fprint(w, dat)
		return
	}

	// We secure passwords with bcrypt with salt (no pepper).
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	passwordHash := string(hash)

	newUser := User{
		Username:       username,
		Email:          email,
		Password:       passwordHash,
		PartieEntiere:  0,
		PartieDecimale: 0,
	}

	// This isn't subject to SQL Injection.
	_, err = database.InsertInto("accounts").
		Columns("*").
		Record(newUser).
		Exec()

	if err != nil {
		datb, _ := webbox.ReadFile("www/signup.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox(err.Error()))
		fmt.Fprint(w, dat)
		return
	}

	// The cookie stored is the right one but will be loaded with the SQL Injection if it has one.
	userSession := &cookieUser{
		Username:      newUser.Username,
		Authenticated: true,
	}

	session.Values["user"] = userSession

	err = session.Save(r, w)

	datb, _ := webbox.ReadFile("www/connected.html")
	dat := newUser.formatConnected(string(datb))
	fmt.Fprint(w, dat)
}

func login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err == nil {
		if val, ok := session.Values["user"]; ok {
			cookieU := val.(cookieUser)
			user, err := getUserByUsername(cookieU.Username)
			if err != nil {
				session.Values["user"] = User{}
				session.Options.MaxAge = -1
				err = session.Save(r, w)
				http.Redirect(w, r, "https://ctf0.yewolf.ovh/login", http.StatusSeeOther)
				return
			}
			datb, _ := webbox.ReadFile("www/connected.html")
			dat := user.formatConnected(string(datb))
			fmt.Fprint(w, dat)
			return
		}
	}
	if r.ContentLength == 0 {
		datb, _ := webbox.ReadFile("www/login.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", "")
		fmt.Fprint(w, dat)
		return
	}
	username := r.FormValue("username")
	passwd := r.FormValue("password")

	if username == "" {
		datb, _ := webbox.ReadFile("www/login.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox("Username can't be empty."))
		fmt.Fprint(w, dat)
		return
	}

	if passwd == "" {
		datb, _ := webbox.ReadFile("www/login.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox("Password can't be empty."))
		fmt.Fprint(w, dat)
		return
	}

	user, err := getUserByUsername(username)
	if err != nil {
		datb, _ := webbox.ReadFile("www/login.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox(err.Error()))
		fmt.Fprint(w, dat)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd)); err != nil {
		datb, _ := webbox.ReadFile("www/login.html")
		dat := string(datb)
		dat = strings.ReplaceAll(dat, "{Error}", errorBox("Invalid password"))
		fmt.Fprint(w, dat)
		return
	}

	userSession := &cookieUser{
		Username:      user.Username,
		Authenticated: true,
	}

	session.Values["user"] = userSession

	err = session.Save(r, w)

	datb, _ := webbox.ReadFile("www/connected.html")
	dat := user.formatConnected(string(datb))
	fmt.Fprint(w, dat)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["user"] = User{}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.SetCookie(w, &http.Cookie{Name: "cookie-name"})
	http.Redirect(w, r, "https://ctf0.yewolf.ovh/login", http.StatusSeeOther)
}

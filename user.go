package main

import (
	"encoding/json"
	"math"
)

type User struct {
	Username       string `db:"username"`
	Email          string `db:"email"`
	Password       string `db:"password"` // This is stored as a hash.
	PartieEntiere  int64  `db:"entiere"`
	PartieDecimale int    `db:"decimale"`
}

func getUserByEmail(email string) (user *User, err error) {
	// SQL Injection is not possible here.
	err = database.Select("*").
		From("accounts").
		Where("email=$1", email).
		QueryStruct(user)
	return
}

func getUserByUsername(username string) (user *User, err error) {
	// SQL Injection is possible here.
	user = &User{}
	err = database.SQL(`select * from "accounts" where username='` + username + `'`).QueryStruct(user)
	return
}

func emailExist(email string) bool {
	var u []User
	d, _ := database.Select("*").
		From("accounts").
		Where("email=$1", email).
		QueryJSON()
	e := json.Unmarshal(d, &u)
	if e != nil {
		// The request can be valid but answer empty so we check for this.
		return false
	}
	return true
}

func usernameExist(username string) bool {
	var u []User
	d, _ := database.Select("*").
		From("accounts").
		Where("username=$1", username).
		QueryJSON()

	e := json.Unmarshal(d, &u)
	if e != nil {
		// Same as l40.
		return false
	}
	return true
}

func (u *User) save() {
	// We only save the money stored in the account.
	database.Update("accounts").
		Set("entiere", u.PartieEntiere).
		Set("decimale", u.PartieDecimale).
		Where("username=$1", u.Username).
		Exec()
}

/*
	You should test those two by yourself if you wanna be sure, but they work and avoid going into negatives.
*/

func (u *User) add(addPartieEntiere int64, addPartieDecimale int) (E int64, D int) {
	E = u.PartieEntiere + addPartieEntiere + int64(math.Floor(float64(u.PartieDecimale+addPartieDecimale)/100.0))
	D = (u.PartieDecimale + addPartieDecimale) % 100
	return
}

func (u *User) sub(subPartieEntiere int64, subPartieDecimale int) (E int64, D int) {
	E = u.PartieEntiere - subPartieEntiere + int64(math.Floor(float64(u.PartieDecimale-subPartieDecimale)/100.0))
	D = (u.PartieDecimale - subPartieDecimale) % 100
	if D < 0 {
		D = 100 + D
	}
	return
}

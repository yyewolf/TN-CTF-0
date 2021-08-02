package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Save some spaces for the error formating (it was specifically made for strings or errors but it would work with anything).
func errorBox(err interface{}) string {
	str := fmt.Sprint(err)
	return `<div class="alert alert-danger" role="alert">
				` + str + `
			</div>`
}

func textAmountToVal(amount string) (E int64, D int, err error) {
	// Has to be in this format : 'XX.XX'
	splt := strings.Split(amount, ".")
	if len(splt) != 2 {
		err = errors.New("Wrong format for amount.")
		return
	}
	E, err = strconv.ParseInt(splt[0], 10, 64)
	if err != nil {
		return
	}
	// We don't want people to be transfering negative amounts.
	if E < 0 {
		E = 0
	}
	D, err = strconv.Atoi(splt[1])
	if err != nil {
		return
	}
	if D < 0 {
		D = 0
	}
	return
}

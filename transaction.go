package main

import (
	"errors"
)

func refillAdmin() {
	// We refill the admin accounts so that anyone can transfer money to themselves.
	database.Exec("update accounts set entiere=10000000000 where username='admin' and entiere <= 1000")
}

// Used to transfer money from one to another
func (u *User) transfer(to string, amountE int64, amountD int) (err error) {
	toUser, err := getUserByUsername(to)
	if err != nil {
		return err
	}

	// Check for negative substraction so that we avoid going in debt.
	restantE, restantD := u.sub(amountE, amountD)
	if restantE < 0 {
		return errors.New("You can't go in debt.")
	}

	u.PartieEntiere = restantE
	u.PartieDecimale = restantD
	toUser.PartieEntiere, toUser.PartieDecimale = toUser.add(amountE, amountD)
	u.save()
	toUser.save()

	refillAdmin()
	return nil
}

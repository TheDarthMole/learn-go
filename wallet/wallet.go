package wallet

import (
	"errors"
	"fmt"
)

type Wallet struct {
	balance Bitcoin
}

type Bitcoin int

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

var ErrNegFunds = errors.New("you cannot withdraw negative funds")
var ErrInsFunds = errors.New("you have insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount < 0 {
		return ErrNegFunds
	}

	if w.Balance() < amount {
		return ErrInsFunds
	}

	w.balance -= amount
	return nil
}

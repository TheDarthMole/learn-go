package bank

import "reflect"

type Person struct {
	Name    string
	Address Address
	NINo    string
}

type Address struct {
	HouseNameOrNum string
	Street         string
	County         string
	Country        string
	Postcode       string
}

type Account struct {
	Owners         []Person
	InitialBalance float64
	Contacts       []Account
}

type Transaction struct {
	From *Account
	To   *Account
	Sum  float64
}

type Ledger struct {
	Ledger []Transaction
}

var transactions = Ledger{}

func Reduce[A, B any](collection []A, accumulator func(B, A) B, initialValue B) (retVal B) {
	retVal = initialValue
	for _, val := range collection {
		retVal = accumulator(retVal, val)
	}
	return
}

func Find[A any](items []A, predicate func(A) bool) (value A, found bool) {
	for _, v := range items {
		if predicate(v) {
			return v, true
		}
	}
	return
}

func (l *Ledger) reset() {
	l.Ledger = make([]Transaction, 0)
}

func (a *Account) Balance() float64 {

	accumulator := func(currentBalance float64, test Transaction) float64 {
		if test.To == a {
			currentBalance += test.Sum
		}
		if test.From == a {
			currentBalance -= test.Sum
		}
		return currentBalance
	}
	return Reduce(transactions.Ledger, accumulator, a.InitialBalance)
}

func (a *Account) IsJointAccount() bool {
	return len(a.Owners) > 1
}

func (a *Account) Send(to *Account, value float64) bool {
	if value > a.Balance() || value < 0 {
		return false
	}
	trans := Transaction{
		From: a,
		To:   to,
		Sum:  value,
	}
	transactions.Ledger = append(transactions.Ledger, trans)
	return true
}

func (a *Account) IsContact(contact Account) bool {
	find := func(p Account) bool {
		return reflect.DeepEqual(p, contact)
	}

	_, found := Find(a.Contacts, find)

	return found
}

func (a *Account) AddContact(contact *Account) {
	if !a.IsContact(*contact) {
		a.Contacts = append(a.Contacts, *contact)
	}
}

func (a *Account) resetContacts() {
	a.Contacts = make([]Account, 0)
}

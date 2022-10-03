package bank

import (
	"errors"
	"reflect"
)

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
	owners         []Person
	initialBalance float64
	contacts       []Account
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

func Filter[A any](items []A, predicate func(A) bool) (filtered []A) {
	for _, v := range items {
		if predicate(v) {
			filtered = append(filtered, v)
		}
	}
	return
}

func NewAccount(owner []Person, initialBalance float64) Account {
	return Account{
		owners:         owner,
		initialBalance: initialBalance,
		contacts:       nil,
	}
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
	return Reduce(transactions.Ledger, accumulator, a.initialBalance)
}

func (a *Account) IsJointAccount() bool {
	return len(a.owners) > 1
}

func (a *Account) Send(to *Account, value float64) error {
	if value > a.Balance() {
		return errors.New("insufficient funds")
	}
	if value < 0 {
		return errors.New("can not send negative balance")
	}

	trans := Transaction{
		From: a,
		To:   to,
		Sum:  value,
	}
	transactions.Ledger = append(transactions.Ledger, trans)
	return nil
}

func (a *Account) IsContact(contact Account) bool {
	find := func(p Account) bool {
		return reflect.DeepEqual(p, contact)
	}

	_, found := Find(a.contacts, find)

	return found
}

func (a *Account) AddContact(contact *Account) {
	if !a.IsContact(*contact) {
		a.contacts = append(a.contacts, *contact)
	}
}

func (a *Account) RemoveContact(contact Account) error {
	oldContacts := a.Contacts()
	newContacts := Filter(a.contacts, func(p Account) bool {
		return !reflect.DeepEqual(p, contact)
	})

	if len(oldContacts) == len(newContacts) {
		return errors.New("the specified contact was not found")
	}

	a.contacts = newContacts
	return nil
}

func (a *Account) resetContacts() {
	a.contacts = make([]Account, 0)
}

func (a *Account) Contacts() []Account {
	return a.contacts
}

func (a *Account) Owners() []Person {
	return a.owners
}

func (a *Account) InitialBalance() float64 {
	return a.initialBalance
}

func (a *Account) AddOwner(newOwner Person) {
	a.owners = append(a.owners, newOwner)
}

func (a *Account) RemoveOwner(remOwner Person) error {
	oldOwners := a.Owners()
	newOwners := Filter(a.owners, func(p Person) bool {
		return !reflect.DeepEqual(remOwner, p)
	})

	if len(newOwners) == len(oldOwners) {
		return errors.New("the specified owner was not found")
	}
	a.owners = newOwners
	return nil
}

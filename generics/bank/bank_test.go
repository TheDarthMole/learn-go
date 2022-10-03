package bank

import (
	"testing"
)

var (
	nick    = Person{Name: "Nick", Address: Address{}, NINo: "PE1337"}
	luke    = Person{Name: "Luke", Address: Address{}, NINo: "PE101337"}
	nickAcc = NewAccount([]Person{nick}, 500)
	lukeAcc = NewAccount([]Person{luke}, 550)
)

func TestAccount_Send(t *testing.T) {
	// Reset the ledger
	transactions.reset()

	t.Run("test sending money", func(t *testing.T) {
		AssertEqual(t, len(transactions.Ledger), 0, "there should be no transactions yet!")

		want := Transaction{
			From: &nickAcc,
			To:   &lukeAcc,
			Sum:  300,
		}

		if err := nickAcc.Send(&lukeAcc, 300); err != nil {
			t.Errorf("transaction should succeed: %s", err)
		}
		AssertEqual(t, len(transactions.Ledger), 1, "there should be one transaction after making a transaction")
		got := transactions.Ledger[0]
		AssertEqual(t, got, want, "")
	})

	t.Run("cant send negative amounts of money", func(t *testing.T) {
		if err := nickAcc.Send(&lukeAcc, -500); err == nil {
			t.Errorf("should give an error for sending negative balance, got: %s", err)
		}
		AssertEqual(t, len(transactions.Ledger), 1, "a new transaction should not have been made for an invalid request")
	})

	t.Run("sending money to yourself doesn't change value", func(t *testing.T) {
		want := Transaction{
			From: &nickAcc,
			To:   &nickAcc,
			Sum:  30,
		}
		originalBal := nickAcc.Balance()
		if err := nickAcc.Send(&nickAcc, 30); err != nil {
			t.Errorf("should be able to send balance to myself: %s", err)
		}
		AssertEqual(t, len(transactions.Ledger), 2, "should have created a new transaction0")
		AssertEqual(t, transactions.Ledger[1], want, "should have a transaction to/from nickAcc")
		AssertEqual(t, nickAcc.Balance(), originalBal, "balances should be the same")
	})
}

func TestAccount_Balance(t *testing.T) {
	// Reset the ledger
	transactions.reset()

	if err := lukeAcc.Send(&nickAcc, 550); err != nil {
		t.Errorf("expected a successful transfer, got: %s", err)
	}
	AssertEqual(t, nickAcc.Balance(), 1050, "Balances not added correctly")
	AssertEqual(t, lukeAcc.Balance(), 0, "Balances not added correctly")

	if err := nickAcc.Send(&lukeAcc, 50); err != nil {
		t.Errorf("expected a successful transfer, got: %s", err)
	}
	AssertEqual(t, nickAcc.Balance(), 1000, "Balances not added correctly")
	AssertEqual(t, lukeAcc.Balance(), 50, "Balances not added correctly")

	t.Run("cant send more than avaliable funds", func(t *testing.T) {
		if err := lukeAcc.Send(&nickAcc, 5000); err == nil {
			t.Errorf("expected transfer to fail due to insufficient funds, got: %s", err)
		}
		AssertEqual(t, nickAcc.Balance(), 1000, "expected values not to change")
		AssertEqual(t, lukeAcc.Balance(), 50, "expected values not to change")
	})
}

func TestAccount_IsJointAccount(t *testing.T) {
	t.Run("one person account", func(t *testing.T) {
		AssertEqual(t, nickAcc.IsJointAccount(), false, "should not be a joint account with one person")
	})

	t.Run("two person account", func(t *testing.T) {
		person1 := Person{}
		person2 := Person{}

		account := NewAccount([]Person{person1, person2}, 1337)

		AssertEqual(t, account.IsJointAccount(), true, "account has two people, should be joint account")
	})

	t.Run("no person account", func(t *testing.T) {
		account := NewAccount(nil, 1337)
		AssertEqual(t, account.IsJointAccount(), false, "no people own this account")
	})

}

func TestAccount_IsContact(t *testing.T) {
	nickAcc.resetContacts()
	lukeAcc.resetContacts()

	t.Run("is a contact", func(t *testing.T) {
		nickAcc.contacts = []Account{lukeAcc}
		AssertEqual(t, nickAcc.IsContact(lukeAcc), true, "Luke was added as a contact to Nick")
	})

	t.Run("is not a contact", func(t *testing.T) {
		AssertEqual(t, lukeAcc.IsContact(nickAcc), false, "Nick should not be a contact for Luke")
	})
}

func TestAccount_AddContact(t *testing.T) {
	nickAcc.resetContacts()
	lukeAcc.resetContacts()

	t.Run("contact doesn't already exist", func(t *testing.T) {
		AssertEqual(t, nickAcc.IsContact(lukeAcc), false, "Luke should not be in contacts list after resetting contacts")
		nickAcc.AddContact(&lukeAcc)
		AssertEqual(t, nickAcc.IsContact(lukeAcc), true, "After adding as contact, should be in contacts list")
	})
}

func TestAccount_RemoveContact(t *testing.T) {
	nickAcc.resetContacts()
	lukeAcc.resetContacts()

	t.Run("removing from empty contacts", func(t *testing.T) {
		AssertEqual(t, len(nickAcc.Contacts()), 0, "Contacts should be empty")
		if err := nickAcc.RemoveContact(lukeAcc); err == nil {
			t.Errorf("wanted to get an error for removing contact when not exists: %s", err)
		}
		AssertEqual(t, len(nickAcc.Contacts()), 0, "Contacts should be empty")
	})

	t.Run("adding contact and removing", func(t *testing.T) {
		nickAcc.AddContact(&lukeAcc)
		AssertEqual(t, nickAcc.IsContact(lukeAcc), true, "should be a contact")
		if err := nickAcc.RemoveContact(lukeAcc); err != nil {
			t.Errorf("got an error when removing a contact: %s", err)
		}

		AssertEqual(t, nickAcc.IsContact(lukeAcc), false, "luke should have been removed as a contact")
	})
}

func AssertEqual[A comparable](t *testing.T, got A, want A, errMessage string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%+v', want '%+v' but they are not equal: %s", got, want, errMessage)
	}
}

package bank

import (
	"testing"
)

var (
	nick    = Person{Name: "Nick", Address: Address{}, NINo: "PE1337"}
	luke    = Person{Name: "Luke", Address: Address{}, NINo: "PE101337"}
	nickAcc = Account{Owners: []Person{nick}, InitialBalance: 500}
	lukeAcc = Account{Owners: []Person{luke}, InitialBalance: 550}
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

		success := nickAcc.Send(&lukeAcc, 300)
		AssertEqual(t, success, true, "transaction should succeed")
		AssertEqual(t, len(transactions.Ledger), 1, "there should be one transaction after making a transaction")
		got := transactions.Ledger[0]
		AssertEqual(t, got, want, "")
	})

	t.Run("cant send negative amounts of money", func(t *testing.T) {
		success := nickAcc.Send(&lukeAcc, -500)
		AssertEqual(t, success, false, "transaction should fail due to invalid value")
		AssertEqual(t, len(transactions.Ledger), 1, "a new transaction should not have been made for an invalid request")
	})

	t.Run("sending money to yourself doesn't change value", func(t *testing.T) {
		want := Transaction{
			From: &nickAcc,
			To:   &nickAcc,
			Sum:  30,
		}
		originalBal := nickAcc.Balance()
		success := nickAcc.Send(&nickAcc, 30)
		AssertEqual(t, success, true, "should be able to send money to myself")
		AssertEqual(t, len(transactions.Ledger), 2, "should have created a new transaction0")
		AssertEqual(t, transactions.Ledger[1], want, "should have a transaction to/from nickAcc")
		AssertEqual(t, nickAcc.Balance(), originalBal, "balances should be the same")
	})
}

func TestAccount_Balance(t *testing.T) {
	// Reset the ledger
	transactions.reset()

	success := lukeAcc.Send(&nickAcc, 550)
	AssertEqual(t, success, true, "expected a successful transfer")
	AssertEqual(t, nickAcc.Balance(), 1050, "Balances not added correctly")
	AssertEqual(t, lukeAcc.Balance(), 0, "Balances not added correctly")

	success = nickAcc.Send(&lukeAcc, 50)
	AssertEqual(t, success, true, "expected a successful transfer")
	AssertEqual(t, nickAcc.Balance(), 1000, "Balances not added correctly")
	AssertEqual(t, lukeAcc.Balance(), 50, "Balances not added correctly")

	t.Run("cant send more than avaliable funds", func(t *testing.T) {
		success = lukeAcc.Send(&nickAcc, 5000)
		AssertEqual(t, success, false, "expected transfer to fail due to insufficient funds")
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

		account := Account{
			Owners:         []Person{person1, person2},
			InitialBalance: 1337,
		}

		AssertEqual(t, account.IsJointAccount(), true, "account has two people, should be joint account")
	})

	t.Run("no person account", func(t *testing.T) {
		account := Account{
			Owners:         nil,
			InitialBalance: 0,
		}

		AssertEqual(t, account.IsJointAccount(), false, "no people own this account")
	})

}

func TestAccount_IsContact(t *testing.T) {
	nickAcc.resetContacts()
	lukeAcc.resetContacts()

	t.Run("is a contact", func(t *testing.T) {
		nickAcc.Contacts = []Account{lukeAcc}
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

func AssertEqual[A comparable](t *testing.T, got A, want A, errMessage string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%+v', want '%+v' but they are not equal: %s", got, want, errMessage)
	}
}

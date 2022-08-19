package wallet

import (
	"testing"
)

func TestWallet(t *testing.T) {

	assertBalance := func(t testing.TB, wallet Wallet, want Bitcoin) {
		t.Helper()
		got := wallet.Balance()
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

	assertError := func(t testing.TB, got error, want error) {
		t.Helper()
		if got == nil {
			t.Error("expected an error but did not get one")
		}

		if got != want {
			t.Errorf("got '%s', want error '%s'", got, want)
		}
	}

	assertNoError := func(t testing.TB, got error) {
		t.Helper()
		if got != nil {
			t.Error("got an error when not expecting one")
		}
	}

	t.Run("Deposit 10 BTC", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Test String method", func(t *testing.T) {
		wallet := Wallet{10}
		got := wallet.Balance().String()
		want := "10 BTC"

		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("Have 20 BTC Withdraw 10 BTC", func(t *testing.T) {
		wallet := Wallet{20}
		err := wallet.Withdraw(Bitcoin(10))

		assertNoError(t, err)
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Have 5 BTC Withdraw 10 BTC", func(t *testing.T) {
		startingBalance := Bitcoin(5)
		wallet := Wallet{startingBalance}
		err := wallet.Withdraw(10)

		assertError(t, err, ErrInsFunds)
		assertBalance(t, wallet, startingBalance)
	})

	t.Run("Have 5 BTC Withdraw -10 BTC", func(t *testing.T) {
		startingBalance := Bitcoin(5)
		wallet := Wallet{startingBalance}
		err := wallet.Withdraw(Bitcoin(-10))

		assertError(t, err, ErrNegFunds)
		assertBalance(t, wallet, startingBalance)
	})

}

func TestError(t *testing.T) {

}

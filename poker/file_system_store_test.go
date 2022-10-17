package poker

import (
	"fmt"
	"github.com/approvals/go-approval-tests/utils"
	"os"
	"testing"
)

// file_system_store_test.go
func TestFileSystemStore(t *testing.T) {

	t.Run("league from a reader", func(t *testing.T) {
		jsonFile := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`

		database, cleanupFile := createTempFile(t, jsonFile)
		defer cleanupFile()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		got := store.GetLeague()

		want := League{
			{"Chris", 33},
			{"Cleo", 10},
		}

		AssertLeague(t, got, want)
		got = store.GetLeague()
		AssertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		jsonFile := `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`

		database, cleanupFile := createTempFile(t, jsonFile)
		defer cleanupFile()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		got := store.GetPlayerScore("Chris")

		want := 33

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}

	})

	//file_system_store_test.go
	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		err = store.RecordWin("Chris")
		AssertNoError(t, err)

		got := store.GetPlayerScore("Chris")
		want := 34
		AssertScoreEquals(t, got, want)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		err = store.RecordWin("Pepper")
		utils.AssertEqual(t, nil, err, "Expected no error when recording win")

		got := store.GetPlayerScore("Pepper")
		want := 1
		AssertScoreEquals(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		playerStore, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		league := playerStore.GetLeague()
		if len(league) != 0 {
			t.Errorf("expected empty league")
		}

		err = playerStore.RecordWin("Jenny")
		AssertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)

		got := store.GetLeague()

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		AssertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		AssertLeague(t, got, want)
	})

	t.Run("store wins then open with same score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)
		player1 := "Nick"
		player2 := "Luke"

		AssertNoError(t, store.RecordWin(player1))
		AssertNoError(t, store.RecordWin(player1))
		AssertNoError(t, store.RecordWin(player1))
		AssertNoError(t, store.RecordWin(player2))

		newStore, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		got := newStore.GetLeague()
		want := League{
			Player{
				Name: "Nick",
				Wins: 3,
			},
			Player{
				Name: "Luke",
				Wins: 1,
			},
		}

		AssertLeague(t, got, want)
	})
}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	_, err = tempFile.Write([]byte(initialData))
	if err != nil {
		t.Errorf("got err writing to tempFile when not expecting an error: %q", err)
	}

	removeFile := func() {
		if err := tempFile.Close(); err != nil {
			fmt.Println("error closing file")
		}

		if err = os.Remove(tempFile.Name()); err != nil {
			fmt.Printf("error removing file %s\n", tempFile.Name())
		}
	}

	return tempFile, removeFile
}

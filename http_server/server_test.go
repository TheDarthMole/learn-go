package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGetPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		[]string{},
	}
	server := NewPlayerServer(&store)
	t.Run("returns Pepper's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Unknown", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

// server_test.go
func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		[]string{},
	}
	server := NewPlayerServer(&store)
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{map[string]int{}, []string{}}
	server := NewPlayerServer(&store)

	t.Run("it records wins when POST", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winnder got %q want %q", store.winCalls[0], player)
		}
	})
}

func TestPreventConcurrencyErrors(t *testing.T) {
	store := StubPlayerStore{map[string]int{}, []string{}}
	server := NewPlayerServer(&store)
	iterations := 1000

	player := "Pepper"
	writeRequest := newPostWinRequest(player)
	readRequest := newGetScoreRequest(player)
	genericResponse := httptest.NewRecorder()

	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			server.ServeHTTP(genericResponse, writeRequest)
			server.ServeHTTP(genericResponse, readRequest)
			wg.Done()
		}()
	}

	wg.Wait()

	if len(store.winCalls) != iterations {
		t.Errorf("got error when multiple concurrent writes, expected %d got %d", iterations, len(store.winCalls))
	}
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get the correct status, got %d want %d", got, want)
	}
}

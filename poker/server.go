package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const jsonContentType = "application/json"

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string) error
	GetLeague() League
}

type PlayerServer struct {
	store PlayerStore
	mu    sync.Mutex
	http.Handler
}

type Player struct {
	Name string
	Wins int
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.store = store
	p.mu = sync.Mutex{}
	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	p.Handler = router
	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", jsonContentType)

	if err := json.NewEncoder(w).Encode(p.getLeagueTable()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) getLeagueTable() League {
	return p.store.GetLeague()
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodGet:
		p.showScore(w, player)
	case http.MethodPost:
		p.processWin(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	_, err := fmt.Fprint(w, score)
	if err != nil {
		return
	}
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	err := p.store.RecordWin(player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusAccepted)
}

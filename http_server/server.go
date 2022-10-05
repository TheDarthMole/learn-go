package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

type PlayerServer struct {
	store PlayerStore
	mu    sync.Mutex
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	return &PlayerServer{store, sync.Mutex{}}
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

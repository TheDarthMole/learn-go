package poker

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
)

//go:embed "templates/*"
var templates embed.FS

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
	template *template.Template
}

type Player struct {
	Name string
	Wins int
}

func NewPlayerServer(store PlayerStore) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFS(templates, "templates/game.html")
	if err != nil {
		return nil, err
	}

	p.store = store
	p.template = tmpl
	p.mu = sync.Mutex{}
	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.gameHandler))
	router.Handle("/ws", http.HandlerFunc(p.websocketHandler))
	p.Handler = router
	return p, nil
}

func (p *PlayerServer) websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, _ := upgrader.Upgrade(w, r, nil)
	_, winnerMsg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err.Error())
	}
	message := string(winnerMsg)
	err = p.store.RecordWin(message)
	log.Println(p.store.GetLeague())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (p *PlayerServer) gameHandler(w http.ResponseWriter, _ *http.Request) {

	err := p.template.Execute(w, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error executing template %s", err.Error()), http.StatusInternalServerError)
	}
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

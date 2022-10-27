package poker

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"net/http"
	"strconv"
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
	http.Handler
	template *template.Template
	store    PlayerStore
	game     Game
	mu       sync.Mutex
}

type Player struct {
	Name string
	Wins int
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFS(templates, "templates/game.html")
	if err != nil {
		return nil, err
	}

	p.game = game
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
	conn, _ := wsUpgrader.Upgrade(w, r, nil)

	_, numberOfPlayersMsg, _ := conn.ReadMessage()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))

	// TODO Don't discard the blinds messages!
	p.game.Start(numberOfPlayers, io.Discard)

	_, winnerMsg, _ := conn.ReadMessage()
	p.game.Finish(string(winnerMsg))

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

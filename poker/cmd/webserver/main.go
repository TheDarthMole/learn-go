package main

import (
	"learn-go/poker"
	"log"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {
	store, closeFile, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	defer closeFile()

	if err != nil {
		log.Fatalf("error creating file system player store from %s", dbFileName)
	}

	server := poker.NewPlayerServer(store)

	if err = http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

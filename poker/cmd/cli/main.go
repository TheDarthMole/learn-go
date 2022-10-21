package main

import (
	"fmt"
	"learn-go/poker"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {

	store, closeFile, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	defer closeFile()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Let's play poker")
	fmt.Println("Type '{Name} wins' to record a win")

	game := poker.NewCLI(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter))

	game.PlayPoker()

}

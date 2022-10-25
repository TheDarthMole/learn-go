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

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.StdOutAlerter), store)
	cli := poker.NewCLI(os.Stdin, os.Stdout, game)

	fmt.Println("Let's play poker!")
	fmt.Println("Type '{Name} wins' to record a win")

	cli.PlayPoker()

}

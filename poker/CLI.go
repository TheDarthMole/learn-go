package poker

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func (cli *CLI) PlayPoker() {

	reader := bufio.NewScanner(cli.in)
	fmt.Printf("Enter winner > ")
	reader.Scan()

	if reader.Text() == "exit" {
		fmt.Printf("Thanks for playing!")
		os.Exit(0)
	}
	err := cli.playerStore.RecordWin(extractWinner(reader.Text()))
	if err != nil {
		log.Fatalf("error recording win, %s", err)
	}

	fmt.Println(cli.playerStore.GetLeague())

}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          in,
	}

}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

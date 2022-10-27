package poker

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

const (
	PlayerPrompt         = "Please enter the number of players: "
	BadPlayerInputErrMsg = "Incorrect input, please try again"
)

func (cli *CLI) PlayPoker() {
	_, err := fmt.Fprint(cli.out, PlayerPrompt)
	if err != nil {
		log.Fatalf("failed to log to client output")
	}

	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, err := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	if err != nil {
		_, err := fmt.Fprint(cli.out, BadPlayerInputErrMsg)
		if err != nil {
			log.Fatal("failed to log to client out")
		}
		return
	}

	cli.game.Start(numberOfPlayers, cli.out)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

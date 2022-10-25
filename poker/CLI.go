package poker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func (cli *CLI) getPlayerCount() int {
	fmt.Fprint(cli.out, PlayerPrompt)
	noOfPlayers, _ := strconv.Atoi(cli.readLine())
	//if err != nil {
	//	fmt.Fprint(cli.out, ErrorPlayerPrompt)
	//	return cli.getPlayerCount()
	//}
	return noOfPlayers
}

func (cli *CLI) quit() {
	fmt.Fprint(cli.out, "Thanks for playing!")
	os.Exit(0)
}

func (cli *CLI) scheduleBlindAlerts(noOfPlayers int) {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		fmt.Println(blind)
		//cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += 10 * time.Minute
	}
}

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
	ErrorPlayerPrompt    = "You did not enter a valid integer, try again"
	WinnerPrompt         = "Please enter the winner: "
	BadPlayerInputErrMsg = "Incorrect input, please try again"
)

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, err := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	if err != nil {
		fmt.Fprint(cli.out, BadPlayerInputErrMsg)
		return
	}

	cli.game.Start(numberOfPlayers)

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

package poker

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type CLI struct {
	playerStore PlayerStore
	in          io.Reader
	alerter     BlindAlerter
}

func (cli *CLI) PlayPoker() {
	cli.scheduleBlindAlerts()

	userInput := cli.readLine()

	switch userInput {
	case "exit":
		cli.quit()
	case "quit":
		cli.quit()
	default:
		err := cli.playerStore.RecordWin(extractWinner(userInput))
		if err != nil {
			log.Fatalf("error recording win, %s", err)
		}
		fmt.Println(cli.playerStore.GetLeague())
	}

}

func (cli *CLI) quit() {
	fmt.Printf("Thanks for playing!")
	os.Exit(0)
}

func (cli *CLI) readLine() string {
	reader := bufio.NewScanner(cli.in)
	fmt.Printf("#> ")
	reader.Scan()
	return reader.Text()
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += 10 * time.Minute
	}
}

func NewCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          in,
		alerter:     alerter,
	}

}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

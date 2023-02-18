package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/alock/lotto-alert/config"
	"github.com/alock/lotto-alert/email"
	"github.com/alock/lotto-alert/scrape"
	"github.com/alock/lotto-alert/util"
)

var (
	rootFlagSet  = flag.NewFlagSet("lotto-alert", flag.ExitOnError)
	test         = rootFlagSet.Bool("test", false, "run commands with fake test data")
	dumpFlagSet  = flag.NewFlagSet("dump", flag.ExitOnError)
	winners      = dumpFlagSet.Bool("winners", false, "show config winners only")
	todayFlagSet = flag.NewFlagSet("today", flag.ExitOnError)
	send         = todayFlagSet.String("send", "", "path for .env file to get app specific password and send email")
)

var (
	dumpCmd = &ffcli.Command{
		Name:       "dump",
		ShortUsage: "lotto-alert [-test] dump [-winners]",
		FlagSet:    dumpFlagSet,
		ShortHelp:  "output PA lotto results from 1/1/23 till today",
		Exec: func(ctx context.Context, args []string) error {
			winningNumbersMap, sortedDates := scrape.GetWinningNumbers(*test)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Date", "Winning Number", "Prize Value", "Prize Reason", "Winner?"})
			table.SetBorder(false)
			for _, date := range sortedDates {
				prize := config.GetDatesPrizeInfo(date)
				winningNum := winningNumbersMap[date]
				if !*winners || (*winners && config.EmailStruct.Participants[winningNum] != "") {
					table.Append([]string{
						util.GetStringOfDate(date),
						util.PadLottoInt(winningNum),
						strconv.Itoa(prize.Amount),
						prize.Reason,
						config.EmailStruct.Participants[winningNum],
					})
				}
			}
			table.Render()
			return nil
		},
	}

	todayCmd = &ffcli.Command{
		Name:       "today",
		ShortUsage: "lotto-alert [-test] today [-send]",
		FlagSet:    todayFlagSet,
		ShortHelp:  "output PA lotto result for today only",
		Exec: func(ctx context.Context, args []string) error {
			exactTime := time.Now()
			today := util.TruncateToDayValue(exactTime)
			winningNumbersMap, _ := scrape.GetWinningNumbers(*test)
			todaysNumber := winningNumbersMap[today]
			if todaysNumber == 0 {
				fmt.Printf("winning was 0 at %v re-run later\n", exactTime)
				return nil
			}
			prize := config.GetDatesPrizeInfo(today)
			fmt.Println("date:", util.GetStringOfDate(today))
			fmt.Println("winning number:", util.PadLottoInt(todaysNumber))
			fmt.Println("prize amount:", prize.Amount)
			fmt.Println("reason:", prize.Reason)
			emailAddress, ok := config.EmailStruct.Participants[todaysNumber]
			if ok {
				fmt.Println("outcome: we have a winner")
				message := email.GetMessage(today, todaysNumber, prize)
				fmt.Println("email:", emailAddress)
				fmt.Println("message:", message)
				if *send != "" {
					return email.SendEmail(emailAddress, message, *send)
				}
				return nil
			}
			fmt.Println("outcome: no winner today")
			return nil
		},
	}
)

func main() {
	config.LoadConfigs()
	root := &ffcli.Command{
		ShortUsage:  "lotto-alert [-test] <subcommand>",
		Subcommands: []*ffcli.Command{dumpCmd, todayCmd},
		FlagSet:     rootFlagSet,
		Exec:        func(context.Context, []string) error { return flag.ErrHelp },
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/alock/lotto-alert/config"
	"github.com/alock/lotto-alert/prize"
	"github.com/alock/lotto-alert/scrape"
)

var winners bool

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "output PA lotto results from 2023 till today",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentYear := time.Now().Year()
		years := make([]int, 0, currentYear-2023+1)
		for y := 2023; y <= currentYear; y++ {
			years = append(years, y)
		}
		winningNumbersMap, sortedDates := scrape.GetWinningNumbers(testMode, years...)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "Winning Number", "Prize Value", "Prize Reason", "Winner?"})
		table.SetBorder(false)
		for _, date := range sortedDates {
			p := prize.ForDate(date)
			winningNum := winningNumbersMap[date]
			if !winners || (winners && config.EmailStruct.Participants[winningNum] != "") {
				table.Append([]string{
					date.Format("01/02/06"),
					fmt.Sprintf("%03d", winningNum),
					strconv.Itoa(p.Amount),
					p.Reason,
					config.EmailStruct.Participants[winningNum],
				})
			}
		}
		table.Render()
		return nil
	},
}

func init() {
	dumpCmd.Flags().BoolVar(&winners, "winners", false, "show config winners only")
	rootCmd.AddCommand(dumpCmd)
}

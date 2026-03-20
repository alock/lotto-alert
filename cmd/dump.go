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

const startYear = 2023

var (
	winners bool
	all     bool
	year    int
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "output PA lotto results (default: current year)",
	Long: `Output PA lotto results in a table.

By default shows the current year only. Use --all for all years
since 2023, or --year to pick a specific year.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentYear := time.Now().Year()

		var years []int
		switch {
		case all:
			for y := startYear; y <= currentYear; y++ {
				years = append(years, y)
			}
		case cmd.Flags().Changed("year"):
			if year < startYear || year > currentYear {
				return fmt.Errorf("year must be between %d and %d", startYear, currentYear)
			}
			years = []int{year}
		default:
			years = []int{currentYear}
		}

		winningNumbersMap, sortedDates, err := scrape.GetWinningNumbers(testMode, years...)
		if err != nil {
			return err
		}
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
	dumpCmd.Flags().BoolVar(&all, "all", false, "show all years since 2023")
	dumpCmd.Flags().IntVar(&year, "year", 0, "show a specific year")
	dumpCmd.MarkFlagsMutuallyExclusive("all", "year")
	rootCmd.AddCommand(dumpCmd)
}

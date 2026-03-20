package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/alock/lotto-alert/config"
	"github.com/alock/lotto-alert/holiday"
)

var cfgYear int

var dumpCfgCmd = &cobra.Command{
	Use:   "dump-cfg",
	Short: "output baked config (participants and computed holidays)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("=== Participants ===")
		for num, email := range config.EmailStruct.Participants {
			fmt.Printf("%03d -> %s\n", num, email)
		}
		fmt.Printf("\n=== Holidays for %d ===\n", cfgYear)
		holidays := holiday.ForYear(cfgYear)

		dates := make([]time.Time, 0, len(holidays))
		for d := range holidays {
			dates = append(dates, d)
		}
		sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })

		for _, d := range dates {
			h := holidays[d]
			fmt.Printf("%s  $%-4d  %s\n", d.Format("01/02"), h.Amount, h.Name)
		}
		return nil
	},
}

func init() {
	dumpCfgCmd.Flags().IntVar(&cfgYear, "year", time.Now().Year(), "year to display holidays for")
	rootCmd.AddCommand(dumpCfgCmd)
}

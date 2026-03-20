package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/alock/lotto-alert/config"
	"github.com/alock/lotto-alert/email"
	"github.com/alock/lotto-alert/prize"
	"github.com/alock/lotto-alert/scrape"
)

var send string

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "output PA lotto result for today only",
	RunE: func(cmd *cobra.Command, args []string) error {
		exactTime := time.Now()
		today := time.Date(exactTime.Year(), exactTime.Month(), exactTime.Day(), 0, 0, 0, 0, exactTime.Location())
		winningNumbersMap, _ := scrape.GetWinningNumbers(testMode, exactTime.Year())
		todaysNumber := winningNumbersMap[today]
		if todaysNumber == 0 {
			fmt.Printf("winning number was 0 at %v re-run later after 16:00:00 PST\n", exactTime.Format(time.DateTime))
			return nil
		}
		p := prize.ForDate(today)
		fmt.Println("date:", today.Format("01/02/06"))
		fmt.Println("winning number:", fmt.Sprintf("%03d", todaysNumber))
		fmt.Println("prize amount:", p.Amount)
		fmt.Println("reason:", p.Reason)
		emailAddress, ok := config.EmailStruct.Participants[todaysNumber]
		if ok {
			fmt.Println("outcome: we have a winner")
			message := email.GetMessage(today, todaysNumber, p)
			fmt.Println("email:", emailAddress)
			fmt.Println("message:", message)
			if send != "" {
				return email.SendEmail(emailAddress, message, send)
			}
			return nil
		}
		fmt.Println("outcome: no winner today")
		return nil
	},
}

func init() {
	todayCmd.Flags().StringVar(&send, "send", "", "path for .env file to get app specific password and send email")
	rootCmd.AddCommand(todayCmd)
}

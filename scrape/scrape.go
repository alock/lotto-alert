package scrape

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const webpageFormat = "https://www.palottery.state.pa.us/Games/Print-Past-Winning-Numbers.aspx?id=28&year=%d&print=1"

func getFakeNums() (map[time.Time]int, []time.Time) {
	nyd := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)
	randomDay := time.Date(time.Now().Year(), 2, 8, 0, 0, 0, 0, time.Local)
	presDay := time.Date(time.Now().Year(), 2, 20, 0, 0, 0, 0, time.Local)
	mommaDay := time.Date(time.Now().Year(), 5, 14, 0, 0, 0, 0, time.Local)
	today := truncateToDay(time.Now())
	return map[time.Time]int{
			nyd:       770,
			randomDay: 123,
			today:     770,
			presDay:   889,
			mommaDay:  52,
		},
		[]time.Time{
			nyd,
			randomDay,
			today,
			presDay,
			mommaDay,
		}
}

// GetWinningNumbers fetches and parses lottery results for the given years.
// If testData is true, returns fake data for testing.
func GetWinningNumbers(testData bool, years ...int) (map[time.Time]int, []time.Time) {
	if testData {
		return getFakeNums()
	}
	m := make(map[time.Time]int)
	for _, year := range years {
		url := fmt.Sprintf(webpageFormat, year)
		text, err := getHtmlPage(url)
		if err != nil {
			log.Fatal(err)
		}
		data := parsePaLottoResults(text)
		for i := 0; i < len(data); i += 2 {
			d, err := time.ParseInLocation("1/2/2006", data[i], time.Local)
			if err != nil {
				log.Fatal(err)
			}
			num, err := strconv.Atoi(data[i+1])
			if err != nil {
				log.Fatal(err)
			}
			m[truncateToDay(d)] = num
		}
	}

	sortedDates := make([]time.Time, 0, len(m))
	for k := range m {
		sortedDates = append(sortedDates, k)
	}
	sort.Slice(sortedDates, func(i, j int) bool {
		return sortedDates[i].Before(sortedDates[j])
	})

	return m, sortedDates
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func getHtmlPage(webPage string) (string, error) {
	resp, err := http.Get(webPage)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func parsePaLottoResults(text string) (data []string) {
	z := html.NewTokenizer(strings.NewReader(text))
	var content []string
	// while have not hit the </html> tag
	for z.Token().Data != "html" {
		tt := z.Next()
		if tt == html.StartTagToken {
			t := z.Token()
			if t.Data == "td" {
				inner := z.Next()
				if inner == html.TextToken {
					text := (string)(z.Text())
					t := strings.TrimSpace(text)
					if t != "" {
						t = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(t, "\u00a0", ""), "\n", ""), " ", "")
						content = append(content, t)
					}
				}
			}
		}
	}
	return content
}

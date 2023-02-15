package scrape

import (
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/alock/lotto-alert/util"
)

const webpage = "https://www.palottery.state.pa.us/Games/Print-Past-Winning-Numbers.aspx?id=28&year=2023&print=1"

func getFakeNums() (map[time.Time]int, []time.Time) {
	nyd := time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
	randomDay := time.Date(2023, 2, 8, 0, 0, 0, 0, time.Local)
	today := util.TruncateToDayValue(time.Now())
	return map[time.Time]int{
			nyd:       770,
			randomDay: 123,
			today:     770,
		},
		[]time.Time{
			nyd,
			randomDay,
			today,
		}
}

func GetWinningNumbers(testData bool) (map[time.Time]int, []time.Time) {
	// get the winning numbers
	if testData {
		return getFakeNums()
	}
	text, err := getHtmlPage(webpage)
	if err != nil {
		log.Fatal(err)
	}
	data := parsePaLottoResults(text)
	m := make(map[time.Time]int)
	for i := 0; i < len(data); i += 2 {
		// probably don't need to place into a date object but this
		// should have value if the scrapping of data is incorrect
		d, err := time.ParseInLocation("1/2/2006", data[i], time.Local)
		if err != nil {
			log.Fatal(err)
		}
		num, err := strconv.Atoi(data[i+1])
		if err != nil {
			log.Fatal(err)
		}
		m[util.TruncateToDayValue(d)] = num
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

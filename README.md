# lotto-alert

Monitors the PA Pick 3 Evening lottery and sends an email alert when a family member's number hits. Runs daily on a Raspberry Pi via cron.

## How it works

1. Scrapes winning numbers from the PA Lottery website
2. Computes today's prize amount — holidays have boosted prizes (e.g. Christmas = $250, Easter = $100)
3. Checks the winning number against a list of family participants (`config/email.json`)
4. If there's a match, sends an email via iCloud SMTP

Holiday dates are computed in code (no more manual JSON updates each year). The `holiday/` package handles Easter (Computus algorithm), equinoxes/solstices (lookup table), nth-weekday holidays (MLK Day, Thanksgiving, etc.), and collision resolution (when two holidays share a date, the higher prize wins).

## Commands

```
lotto-alert today              # check today's result
lotto-alert today --send .env  # check + send email if winner
lotto-alert dump               # show all results from 2023 to now
lotto-alert dump --winners     # show only winning days
lotto-alert dump-cfg           # show participants and computed holidays
lotto-alert dump-cfg --year 2025  # show holidays for a specific year
```

Add `--test` before any subcommand to use fake data instead of scraping the live site.

## Setup

### Prerequisites
- Go 1.23+
- [iCloud+ custom email domain](https://support.apple.com/en-us/HT212514)
- [App-specific password](https://support.apple.com/en-us/HT204397)

### Configuration

**`config/email.json`** — participants and sender config (see `config/email-template.json` for the format):
```json
{
  "from": "you@icloud.com",
  "fromOverride": "alias@yourdomain.com",
  "participants": {
    "123": "winner@example.com",
    "456": "another@example.com"
  }
}
```

**`.env`** — app-specific password for iCloud SMTP:
```
APP_SPECIFIC=xxxx-xxxx-xxxx-xxxx
```

### Building

```
make build         # build for local arch
make test          # run all tests
make lint          # go vet
make pi            # cross-compile for ARM + scp to pi
make clean         # remove binary
```

Or use `./pi.sh` which just runs `make pi`.

### Raspberry Pi cron

```
20 16 * * * /home/pi/.local/bin/lotto-alert today --send /home/pi/.local/bin/.env
```

Debug cron with:
```
tail -fn 0 /var/log/syslog
```

## Project structure

```
main.go                  # entry point
cmd/                     # cobra commands (today, dump, dump-cfg)
holiday/                 # holiday computation engine
  easter.go              # Computus algorithm
  equinox.go             # equinox/solstice lookup table
  holiday.go             # 37 holiday definitions + collision resolution
  testdata/              # ground truth fixtures from 2023-2026
prize/                   # prize lookup (holiday > Saturday $50 > weekday $30)
scrape/                  # PA Lottery website scraper
email/                   # iCloud SMTP sender
config/                  # embedded participant config (email.json)
```

## Updating for a new year

The holiday engine computes dates automatically — no branch-per-year needed. The only thing that might need a tweak is the equinox/solstice lookup table in `holiday/equinox.go` if the lottery uses non-standard dates. Run `make test` to validate against 4 years of historical data.

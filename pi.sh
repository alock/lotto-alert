#!/bin/sh

env GOOS=linux GOARCH=arm GOARM=7 go build
scp lotto-alert pi.local:/home/pi/.local/bin
rm lotto-alert

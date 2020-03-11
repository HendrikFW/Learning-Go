package main

import (
	"fmt"
	"time"
)

func fmtDate(date time.Time) string {
	return fmt.Sprintf("%02d.%02d.%d %02d:%02d:%02d",
		date.Day(),
		date.Month(),
		date.Year(),
		date.Hour(),
		date.Minute(),
		date.Second())
}

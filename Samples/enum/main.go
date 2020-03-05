package main

import "fmt"

type Weekday int

const (
	Sunday    Weekday = 0
	Monday    Weekday = 1
	Tuesday   Weekday = 2
	Wednesday Weekday = 3
	Thursday  Weekday = 4
	Friday    Weekday = 5
	Saturday  Weekday = 6
)

func (day Weekday) String() string {
	names := [...]string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursaday",
		"Friday",
		"Saturday",
	}

	if day < Sunday || day > Saturday {
		return "Unknown"
	}

	return names[day]
}

func main() {
	fmt.Println("Learning Go: enums")

	var bestDayOfWeek Weekday = Friday

	fmt.Printf("The best day of the week is %s\n", bestDayOfWeek)
	fmt.Printf("%d\n", bestDayOfWeek)
}

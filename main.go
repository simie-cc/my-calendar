package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Event struct {
	EventTime time.Time
	EventDesc string
}

func main() {
	fmt.Println("test")
	events := readAndParse()
	t := template.New("Template_1")
	t.Parse(icsTemplate)

	output, err := os.Create("./docs/cal.ics")
	if err != nil {
		fmt.Println("Cannot create output file: ", err)
		os.Exit(1)
	}
	defer output.Close()

	err = t.Execute(output, map[string]interface{}{
		"today":  time.Now().In(time.UTC).Format("20060102T150405Z"),
		"events": events,
		"formatUTC": func(d time.Time) string {
			return d.In(time.UTC).Format("20060102T150405")
		},
		"format": func(d time.Time) string {
			return d.Format("20060102T150405")
		},
	})
	if err != nil {
		fmt.Println("error template: ", err)
		os.Exit(1)
	}

}

func readAndParse() []Event {
	cal, err := os.Open("calendar.txt")
	if err != nil {
		fmt.Println("Open error: ", err)
		os.Exit(1)
	}
	defer cal.Close()

	lineRegex := regexp.MustCompile("^([[:digit:]]{4})/([[:digit:]]{1,2})/([[:digit:]]{1,2}) (.*)$")
	now := time.Now()
	local := now.Location()

	var events []Event
	scn := bufio.NewScanner(cal)
	for scn.Scan() {
		line := scn.Text()

		sub := lineRegex.FindStringSubmatch(line)
		year, err := strconv.Atoi(sub[1])
		if err != nil {
			fmt.Printf("Date(year) error for line [%s]: %v\n", line, err)
			continue
		}

		month, err := strconv.Atoi(sub[2])
		if err != nil {
			fmt.Printf("Date(month) error for line [%s]: %v\n", line, err)
			continue
		}

		dayOfMonth, err := strconv.Atoi(sub[3])
		if err != nil {
			fmt.Printf("Date(day) error for line [%s]: %v\n", line, err)
			continue
		}

		eventTime := time.Date(year, time.Month(month), dayOfMonth, 0, 0, 0, 0, local)
		if err != nil {
			fmt.Printf("Date error for line [%s]: %v\n", line, err)
			continue
		}

		eventDesc := sub[4]

		fmt.Println("Line: ", line, "===", eventTime.Format("2006/01/02"), eventDesc)
		events = append(events, Event{
			EventTime: eventTime,
			EventDesc: eventDesc,
		})
	}

	return events
}

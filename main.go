package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	EventTime    time.Time
	EventTimeEnd time.Time
	EventDesc    string
}

func main() {

	events := readAndParse()
	temp := template.New("Template_1")
	temp.Parse(icsTemplate)

	output, err := os.Create("./docs/cal.ics")
	if err != nil {
		fmt.Println("Cannot create output file: ", err)
		os.Exit(1)
	}
	defer output.Close()

	now := time.Now()
	data := map[string]interface{}{
		"now":    now.In(time.UTC).Format("20060102T150405Z"),
		"events": events,
		"formatUTC": func(d time.Time) string {
			return d.In(time.UTC).Format("20060102T150405")
		},
		"format": func(d time.Time) string {
			return d.Format("20060102")
		},
	}
	err = generateIcs(temp, now, output, data)
	if err != nil {
		fmt.Println("error template: ", err)
		os.Exit(1)
	}
}

// func shaString(raw string) string {
// 	s := sha512.New()
// 	s.Write([]byte(raw))
// 	return base64.URLEncoding.EncodeToString(s.Sum(nil))
// }

func generateIcs(temp *template.Template, now time.Time, output io.Writer, data interface{}) error {
	return temp.Execute(output, data)
}

func readAndParse() []Event {
	cal, err := os.Open("calendar.txt")
	if err != nil {
		fmt.Println("Open error: ", err)
		os.Exit(1)
	}
	defer cal.Close()

	lineRegex := regexp.MustCompile("^([[:digit:]]{4})/([[:digit:]]{1,2})/([[:digit:]]{1,2})(-([[:digit:]]{4})/([[:digit:]]{1,2})/([[:digit:]]{1,2}))? (.*)$")

	var events []Event
	scn := bufio.NewScanner(cal)
	for scn.Scan() {
		line := scn.Text()

		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		sub := lineRegex.FindStringSubmatch(line)
		eventTime, err := parseEventDate(line, sub, 1)
		if err != nil {
			continue
		}

		var eventTimeEnd time.Time
		if len(sub[4]) > 0 {
			eventTimeEnd, err = parseEventDate(line, sub, 5)
			if err != nil {
				continue
			}
		} else {
			eventTimeEnd = eventTime.AddDate(0, 0, 1)
		}
		eventDesc := sub[8]

		fmt.Println("Line: ", line, "===",
			eventTime.Format("2006/01/02"), eventTimeEnd.Format("2006/01/02"), eventDesc)
		events = append(events, Event{
			EventTime:    eventTime,
			EventTimeEnd: eventTimeEnd,
			EventDesc:    eventDesc,
		})
	}

	return events
}

func parseEventDate(line string, sub []string, start int) (time.Time, error) {
	year, err := strconv.Atoi(sub[start+0])
	if err != nil {
		fmt.Printf("Date(year) error for line [%s]: %v\n", line, err)
		return time.Time{}, err
	}

	month, err := strconv.Atoi(sub[start+1])
	if err != nil {
		fmt.Printf("Date(month) error for line [%s]: %v\n", line, err)
		return time.Time{}, err
	}

	dayOfMonth, err := strconv.Atoi(sub[start+2])
	if err != nil {
		fmt.Printf("Date(day) error for line [%s]: %v\n", line, err)
		return time.Time{}, err
	}

	eventTime := time.Date(year, time.Month(month), dayOfMonth, 0, 0, 0, 0, time.Local)
	if err != nil {
		fmt.Printf("Date error for line [%s]: %v\n", line, err)
		return time.Time{}, err
	}

	return eventTime, nil
}

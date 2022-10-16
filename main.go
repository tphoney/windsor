package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"golang.org/x/net/publicsuffix"
)

var cli struct {
	Username     string   `help:"username" default:"thomas honey"`
	Password     string   `help:"password" default:"password"`
	Participants []string `help:"participants" default:"darragh lewis,ronan okane,gary toal"`
	Day          string   `help:"specify date"`
	Rooms        []string `help:"rooms" default:"28,27"`
	Hour         int      `help:"hour" default:"20"`
}

func main() {
	// parse command line arguments
	kong.Parse(&cli,
		kong.Name("Windsor booking bot"),
		kong.Description("A simple bot to book windsor tennis courts"))
	hour := hoursToSeconds(cli.Hour)
	year, month, day := eightDaysLater()
	if cli.Day != "" {
		day = cli.Day
	}

	fmt.Printf("booking for '%s' with %v in rooms %v @ y:%s m:%s d:%s h:%d:00=%d \n", cli.Username, cli.Participants, cli.Rooms, year, month, day, cli.Hour, hour)
	fmt.Println("starting booking...")
	sleep := time.Duration(1) * time.Second
	// try booking every 1 second for 2 minutes
	for int := 0; int < 120; int++ {
		for _, room := range cli.Rooms {
			booked := loginAndBook(cli.Username, cli.Password, year, month, day, fmt.Sprint(hour), room, cli.Participants)
			if booked {
				fmt.Println("exiting")
				os.Exit(0)
			}
		}
		time.Sleep(sleep)
	}
}

func loginAndBook(username, password, year, month, day, hour, room string, participants []string) (success bool) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, cookieErr := cookiejar.New(&options)
	if cookieErr != nil {
		log.Fatal(cookieErr)
		return false
	}

	client := http.Client{Jar: jar}
	loginURL := fmt.Sprintf("https://www.windsortennis.co.uk/courtbooker/day.php?day=%s&month=%s&year=%s&area=13&room=3", day, month, year)
	resp, postErr := client.PostForm(loginURL,
		url.Values{
			"NewUserPassword": {password},
			"NewUserName":     {username},
			"Action":          {"SetName"},
		})
	if postErr != nil {
		log.Fatal(postErr)
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	contents := buf.String()
	if !strings.Contains(contents, fmt.Sprintf("You are %s", username)) {
		fmt.Println("failed to log in")
		return false
	}
	fmt.Println("logged in")
	// lets book a court
	bookingURL := "http://www.windsortennis.co.uk/courtbooker/edit_entry_handler.php"
	payload := url.Values{
		"name":          {username},
		"description":   {"robot booking"},
		"start_day":     {day},
		"start_month":   {month},
		"start_year":    {year},
		"start_seconds": {hour},
		"end_day":       {day},
		"end_month":     {month},
		"end_year":      {year},
		"end_seconds":   {hour},
		"area":          {"13"}, // this is the dome
		"rooms":         {room},
		"type":          {"A"},
		"create_by":     {username},
		"rep_id":        {"0"},
		"agree":         {"1"},
		"edit_type":     {"series"},
	}
	for i, participant := range participants {
		payload.Add(fmt.Sprintf("participant_%d", i+1), participant)
	}
	resp, postErr = client.PostForm(bookingURL, payload)
	if postErr != nil {
		log.Fatal(postErr)
		return false
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bookedContents := buf.String()
	if strings.Contains(bookedContents, "robot booking") {
		fmt.Println("booked")
		return true
	}
	fmt.Println("failed to book")
	// write the contents to a file
	os.WriteFile("fail.html", []byte(contents), 0644)
	return false
}

// return the date 8 days from now
func eightDaysLater() (year, month, day string) {
	// get current date
	currentDate := time.Now()
	nextWeek := currentDate.AddDate(0, 0, 8)
	return nextWeek.Format("2006"), nextWeek.Format("01"), nextWeek.Format("02")
}

// convert hours to seconds
func hoursToSeconds(hours int) int {
	return hours * 60 * 60
}

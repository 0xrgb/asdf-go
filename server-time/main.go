package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Quiet bool `short:"q" long:"quiet" description:"Do not print error message"`
}

func getTime(url string) (time.Time, error) {
	resp, err := http.Get(url)
	if err != nil {
		return time.Time{}, fmt.Errorf("getTime: request failed: %v", err)
	}

	date := resp.Header["Date"]
	if len(date) != 1 {
		return time.Time{}, errors.New("getTime: wrong date header size")
	}

	serverTime, err := time.Parse(
		time.RFC1123,
		date[0],
	)

	if err != nil {
		return time.Time{}, fmt.Errorf("getTime: cannot parse time %q", date[0])
	}

	localTime := serverTime.Local()
	return localTime, nil
}

func main() {
	urls, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(2)
		} else {
			panic(err)
		}
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	urlColors := []*color.Color{
		color.New(color.FgCyan, color.Underline),
		color.New(color.FgGreen, color.Underline),
		color.New(color.FgMagenta, color.Underline),
		color.New(color.FgCyan, color.Underline),
		color.New(color.FgBlue, color.Underline),
	}

	for _, url := range urls {
		t, err := getTime(url)

		col := rnd.Intn(len(urlColors))
		urlColors[col].Print(url)
		fmt.Print(": ")
		if err != nil {
			if !opts.Quiet {
				color.Red("%v\n", err)
			}
		} else {
			fmt.Printf("%s\n", t.Format("15:04:05"))
		}
	}
}

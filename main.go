package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {

	interval := flag.Int("n", 60, "Interval in seconds between requests")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("USAGE: urlwatcher [-n sec] <URL>")
		os.Exit(1)
	}
	watchurl := flag.Arg(0)

	httpclient := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		// Timeout: 5 * time.Second,
	}

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	for {
		// send http request
		go func() {
			starttime := time.Now()

			req, _ := http.NewRequest(http.MethodGet, watchurl, nil)
			req.Header.Set("User-Agent", "urlwatcher/0.1")

			res, err := httpclient.Do(req)
			if err != nil {
				switch err := err.(type) {
				case *url.Error:
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					if !err.Temporary() {
						os.Exit(1)
					}
				default:
					fmt.Fprintf(os.Stderr, "Error fetching URL: %#v\n", err)
				}

				fmt.Printf("%s\t%d\t%d\t%.2f\n",
					starttime.Format(time.DateTime), -1, -1, time.Since(starttime).Seconds())

				return
			}

			// fmt.Printf("%#v\n", res)

			io.Copy(io.Discard, res.Body)
			res.Body.Close()

			fmt.Printf("%s\t%d\t%d\t%.2f\n",
				starttime.Format(time.DateTime), res.StatusCode, res.ContentLength, time.Since(starttime).Seconds())
		}()

		// wait for the next tick
		<-ticker.C
	}

}

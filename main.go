package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type LogEntry struct {
	time   time.Time
	method string
	url    string
}

var logChannel chan string
var logWg sync.WaitGroup
var httpWg sync.WaitGroup

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var format string
var nginxLogFile string
var logFile string
var prefix string
var ratio int64
var debug bool

func init() {
	flag.StringVar(&format, "format", `$remote_addr [$time_local] "$request" $status $request_length $body_bytes_sent $request_time "$t_size" $read_time $gen_time`, "Log format")
	flag.StringVar(&nginxLogFile, "file", "dummy", "Log file name to read. Read from STDIN if file name is '-'")
	flag.StringVar(&logFile, "log", "-", "File to report timings to, default is stdout")
	flag.StringVar(&prefix, "prefix", "http://localhost", "Url prefix to query")
	flag.Int64Var(&ratio, "ratio", 1, "Replay speed ratio, higher means faster replay speed")
	flag.BoolVar(&debug, "debug", false, "Print extra debugging information")

	logChannel = make(chan string)
}

func mainLoop(reader *NginxReader) {
	var nilTime time.Time
	var lastTime time.Time

	for {
		rec, err := reader.Read()

		if err == io.EOF {
			log.Println("Reached EOF")
			break
		} else {
			checkErr(err)
		}

		if rec.method == "GET" {
			if lastTime != nilTime {
				differenceUnix := rec.time.Sub(lastTime).Nanoseconds()

				if differenceUnix > 0 {
					durationWithRation := time.Duration(differenceUnix / ratio)

					if debug {
						log.Printf("Sleeping for: %.2f seconds", durationWithRation.Seconds())
					}
					time.Sleep(durationWithRation)
				} else {
					if debug {
						log.Println("No need for sleep!")
					}
				}

			}

			lastTime = rec.time

			httpWg.Add(1)
			go fireHttpRequest(rec.url)
		}
	}
}

func fireHttpRequest(url string) {
	defer httpWg.Done()

	path := prefix + url

	if debug {
		log.Println("Querying %s", path)
	}

	startTime := time.Now()
	resp, err := http.Get(path)
	endTime := time.Now()

	if err != nil {
		log.Printf(`ERROR "%s" while querying "%s"`, err, path)
	} else {
		status := resp.StatusCode
		duration := endTime.Sub(startTime).Nanoseconds()
		logMessage := fmt.Sprintf("%d %d %s\n", status, duration, url)

		logChannel <- logMessage
	}
}

func logLoop() {
	defer logWg.Done()

	var writer io.Writer

	switch logFile {
	case "-":
		writer = os.Stdout
	default:
		file, err := os.Create(logFile)
		checkErr(err)
		writer = file
	}

	for logMessage := range logChannel {
		_, err := io.WriteString(writer, logMessage)
		checkErr(err)
	}
}

func main() {
	flag.Parse()

	var logReader io.Reader

	if debug {
		log.Printf("Parsing %s log file\n", nginxLogFile)
		log.Printf("Using format %s", format)
	}

	if nginxLogFile == "dummy" {
		logReader = strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /t/100x100/foo/bar.jpeg HTTP/1.1" 200 1027 2430 0.014 "100x100" 10 1`)
	} else if nginxLogFile == "-" {
		logReader = os.Stdin
	} else {
		file, err := os.Open(nginxLogFile)

		checkErr(err)
		defer file.Close()

		logReader = file
	}

	reader := NewNginxReader(logReader, format)
	log.Println(reader)

	logWg.Add(1)
	go logLoop()

	mainLoop(reader)

	if debug {
		log.Println("Waiting for all http goroutines to stop")
	}

	httpWg.Wait()
	close(logChannel)

	if debug {
		log.Println("Waiting for log goroutine to stop")
	}

	logWg.Wait()
}

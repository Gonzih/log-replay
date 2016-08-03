package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	haProxyTsLayout = "2/Jan/2006:15:04:05.000"
)

// HaproxyReader implements LogReader intefrace
type HaproxyReader struct {
	InputReader  io.Reader
	InputScanner *bufio.Scanner
}

func parseHaproxyTime(timeLocal string) time.Time {
	t, err := time.Parse(haProxyTsLayout, timeLocal)

	checkErr(err)

	return t
}

func parseStringInto(s string, entry *LogEntry) error {
	dateStartI := strings.LastIndex(s, "[") + 1
	dateEndI := strings.LastIndex(s, "]")

	if dateStartI > dateEndI || dateStartI > len(s) || dateEndI > len(s) {
		return fmt.Errorf("Issue with date indexes, start: %d, end: %d, len: %d", dateStartI, dateEndI, len(s))
	}

	requestStartI := strings.Index(s, `"`) + 1
	requestEndI := len(s) - 1

	if requestStartI > requestEndI || requestStartI > len(s) || requestEndI > len(s) {
		return fmt.Errorf("Issue with request indexes, start: %d, end: %d, len: %d", requestStartI, requestEndI, len(s))
	}

	dateString := s[dateStartI:dateEndI]
	requestString := s[requestStartI:requestEndI]

	parsedRequest, err := parseRequest(requestString)

	if err != nil {
		return err
	}

	entry.Method = parsedRequest[0]
	entry.URL = parsedRequest[1]
	entry.Time = parseHaproxyTime(dateString)

	return nil
}

// NewHaproxyReader creates new reader for a haproxy log format using provided io.Reader
func NewHaproxyReader(inputReader io.Reader) LogReader {
	var reader HaproxyReader

	reader.InputReader = inputReader
	reader.InputScanner = bufio.NewScanner(reader.InputReader)

	return &reader
}

func (r *HaproxyReader) Read() (*LogEntry, error) {
	var entry LogEntry

	inputAvailable := r.InputScanner.Scan()

	if inputAvailable {
		parseStringInto(r.InputScanner.Text(), &entry)
	} else {
		return &entry, io.EOF
	}

	err := r.InputScanner.Err()

	if err != nil {
		return &entry, err
	}

	return &entry, nil
}

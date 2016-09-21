package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

const (
	solrProxyTsLayout = "2006-01-02 15:04:05.000"
)

// SolrReader implements LogReader intefrace
type SolrReader struct {
	InputReader  io.Reader
	InputScanner *bufio.Scanner
}

func parseSolrTime(timeLocal string) time.Time {
	t, err := time.Parse(solrProxyTsLayout, timeLocal)

	checkErr(err)

	return t
}

func parseSolrPayload(params string) (string, error) {
	r := regexp.MustCompile(`{(.+)?}`)
	matches := r.FindStringSubmatch(params)
	if len(matches) != 2 {
		return "", fmt.Errorf("Unable to parse solr payload.")
	}
	return matches[1], nil
}

func parseSolrInto(s string, entry *LogEntry) error {
	if len(s) < 23 {
		return fmt.Errorf("This log line does not seem to contain a valid timestamp.")
	}
	dateString := strings.Replace(s[0:23], ",", ".", -1)
	stringParts := strings.Split(s, " ")

	var requestParts [2]string
	for _, part := range stringParts {
		if strings.HasPrefix(part, "path") {
			requestParts[0] = part
		} else if strings.HasPrefix(part, "params") {
			requestParts[1] = part
		}
	}
	//Default solr requests to post to go around query length for GET requests.
	payload, err := parseSolrPayload(requestParts[1])
	if err != nil {
		return err
	}

	path := strings.SplitAfterN(requestParts[0], "=", 2)

	entry.Method = "POST"
	entry.URL = path[1]
	entry.Time = parseSolrTime(dateString)
	entry.Payload = payload
	return nil
}

// NewSolrReader creates new reader for a solr log format using provided io.Reader
func NewSolrReader(inputReader io.Reader) LogReader {
	var reader SolrReader

	reader.InputReader = inputReader
	reader.InputScanner = bufio.NewScanner(reader.InputReader)

	return &reader
}

func (r *SolrReader) Read() (*LogEntry, error) {
	var entry LogEntry

	inputAvailable := r.InputScanner.Scan()

	if inputAvailable {
		parseSolrInto(r.InputScanner.Text(), &entry)
	} else {
		return &entry, io.EOF
	}

	err := r.InputScanner.Err()

	if err != nil {
		return &entry, err
	}

	return &entry, nil
}

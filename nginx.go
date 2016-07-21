package main

import (
	"errors"
	"fmt"
	"github.com/satyrius/gonx"
	"io"
	"strings"
	"time"
)

const (
	nginxTimeLayout = "2/Jan/2006:15:04:05 -0700"
)

type NginxReader struct {
	GonxReader *gonx.Reader
}

func parseRequest(requestString string) ([]string, error) {
	parsedRequest := strings.SplitN(requestString, " ", 3)

	if len(parsedRequest) != 3 {
		return parsedRequest, errors.New(fmt.Sprintf("ERROR while parsing string: %s", requestString))
	}

	return parsedRequest, nil
}

func parseNginxTime(timeLocal string) time.Time {
	t, err := time.Parse(nginxTimeLayout, timeLocal)

	checkErr(err)

	return t
}

func NewNginxReader(inputReader io.Reader, format string) LogReader {
	var reader NginxReader
	reader.GonxReader = gonx.NewReader(inputReader, format)

	return &reader
}

func (r *NginxReader) Read() (*LogEntry, error) {
	var entry LogEntry

	rec, err := r.GonxReader.Read()

	if err != nil {
		return &entry, err
	}

	timeLocal, err := rec.Field("time_local")

	if err != nil {
		return &entry, err
	}

	requestString, err := rec.Field("request")

	if err != nil {
		return &entry, err
	}

	parsedRequest, err := parseRequest(requestString)

	if err != nil {
		return &entry, err
	}

	entry.method = parsedRequest[0]
	entry.url = parsedRequest[1]
	entry.time = parseNginxTime(timeLocal)

	return &entry, nil
}

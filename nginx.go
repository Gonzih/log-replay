package main

import (
	"errors"
	"fmt"
	"github.com/satyrius/gonx"
	"io"
	"strings"
	"time"
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

func NewNginxReader(logReader io.Reader, format string) *NginxReader {
	var reader NginxReader
	reader.GonxReader = gonx.NewReader(logReader, format)

	return &reader
}

func parseTimeLocal(timeLocal string) time.Time {
	layout := "2/Jan/2006:15:04:05 -0700"

	t, err := time.Parse(layout, timeLocal)

	checkErr(err)

	return t
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
	entry.time = parseTimeLocal(timeLocal)

	return &entry, nil
}

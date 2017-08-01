# Replay Nginx/Haproxy/SOLR logs

[![MIT License][license-image]][license-url]
[![Build Status](https://travis-ci.org/Gonzih/log-replay.svg?branch=master)](https://travis-ci.org/Gonzih/log-replay)

## Installation

```
go get -u github.com/Gonzih/log-replay
```

## Usage

```
Usage of log-replay:
  -debug
    	Print extra debugging information
  -file string
    	Log file name to read. Read from STDIN if file name is '-' (default "-")
  -file-type string
    	Input log type (nginx, haproxy or solr) (default "nginx")
  -format string
    	Nginx log format (default "$remote_addr [$time_local] \"$request\" $status $request_length $body_bytes_sent $request_time \"$t_size\" $read_time $gen_time")
  -log string
    	File to report timings to, default is stdout (default "-")
  -prefix string
    	URL prefix to query (default "http://localhost")
  -ratio int
    	Replay speed ratio, higher means faster replay speed (default 1)
  -skip-sleep
    	Skip sleep between http calls based on log timestapms
  -timeout int
    	Request timeout in milliseconds, 0 means no timeout (default 60000)
```

```bash
# Replay access log
log-replay --file my-acces.log --debug --log out.log

# Duplicate traffic on the staging host
tail -f /var/log/acces.log | log-replay --prefix http://staging-host --log staging.log --skip-sleep
```

## Output log format

Log is tab separated values:

```
status	start-time	duration	url payload err

# Examples
200	1469792268	629904766	/my-url
500	1469792268	629904766	/my-url	Get http://localhost/another-url: dial tcp [::1]:80: getsockopt: connection refused
```

* status is integer
* start-time is unix timestamp in seconds
* duration is in nanoseconds
* url is full url with prefix
* payload is stringified post data
* error is go lang error formatted to string and is optional

## Only GET?

Nginx/Haproxy logs are currently limited to GET only.
SOLR requests will use post format for everything, as a way to subvert GET length limitations.

## Log formats

* To correctly use the solr adapter, it is required that the log4 pattern is configured as follows:

```
<PatternLayout>
  <pattern>%d %p %C{1.} [%t] %m%n%ex</pattern>
</PatternLayout>
```

## License

[MIT](LICENSE)

[license-url]: LICENSE

[license-image]: https://img.shields.io/github/license/mashape/apistatus.svg

[capture]: capture.png

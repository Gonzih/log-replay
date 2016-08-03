# Replay Nginx/Haproxy logs

[![MIT License][license-image]][license-url]
[![Build Status](https://travis-ci.org/Gonzih/log-replay.svg?branch=master)](https://travis-ci.org/Gonzih/log-replay)

## Installation

```
go get -u github.com/Gonzih/log-replay
```

## Usage

```
$ log-replay:
  -debug
        Print extra debugging information
  -file string
        Log file name to read. Read from STDIN if file name is '-' (default "-")
  -file-type string
        Input log type (nginx or haproxy) (default "nginx")
  -format string
        Nginx log format (default "$remote_addr [$time_local] \"$request\" $status $request_length $body_bytes_sent $request_time \"$t_size\" $read_time $gen_time")
  -log string
        File to report timings to, default is stdout (default "-")
  -prefix string
        Url prefix to query (default "http://localhost")
  -ratio int
        Replay speed ratio, higher means faster replay speed (default 1)
```

```bash
# Replay access log
log-replay --file my-acces.log --debug --log out.log

# Duplicate traffic on the staging host
tail -f /var/log/acces.log | log-replay --prefix http://staging-host --log staging.log
```

## Output log format

Log is tab separated values:

```
status	start-time	duration	url	err

# Examples
200	1469792268	629904766	/my-url
500	1469792268	629904766	/my-url	Get http://localhost/another-url: dial tcp [::1]:80: getsockopt: connection refused
```

* status is integer
* start-time is unix timestamp in seconds
* duration is in nanoseconds
* url is full url with prefix
* error is go lang error formatted to string and is optional

## Only GET?

Yeah, for now only get requests are replayed, not sure if there is need to replay other http methods.

## License

[MIT](LICENSE)

[license-url]: LICENSE

[license-image]: http://img.shields.io/badge/license-MIT-blue.svg?style=flat

[capture]: capture.png

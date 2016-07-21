# Replay Nginx logs [![Build Status](https://travis-ci.org/Gonzih/replay-nginx-log.svg?branch=master)](https://travis-ci.org/Gonzih/tmuxman)

## Installation

```
go get github.com/Gonzih/replay-nginx-log
```

## Usage

```
$ replay-nginx-log --help

Usage of replay-nginx-log:
  -debug
        Print extra debugging information
  -file string
        Log file name to read. Read from STDIN if file name is '-' (default "-")
  -format string
        Log format (default "$remote_addr [$time_local] \"$request\" $status $request_length $body_bytes_sent $request_time \"$t_size\" $read_time $gen_time")
  -log string
        File to report timings to, default is stdout (default "-")
  -prefix string
        Url prefix to query (default "http://localhost")
  -ratio int
        Replay speed ratio, higher means faster replay speed (default 1)
```

```bash
# Replay access log
replay-nginx-log --file my-acces.log --debug --log out.log

# Duplicate traffic on the staging host
tail -f /var/log/acces.log | replay-nginx-log --prefix http://staging-host --log staging.log
```

## Only GET?

Yeah, for now only get requests are replayed, not sure if there is need to replay other kind of requests.

# Replay Nginx/Haproxy logs [![Build Status](https://travis-ci.org/Gonzih/log-replay.svg?branch=master)](https://travis-ci.org/Gonzih/log-replay)

## Installation

```
go get github.com/Gonzih/log-replay
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

## Only GET?

Yeah, for now only get requests are replayed, not sure if there is need to replay other http methods.

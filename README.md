# Replay Nginx logs

## Installation

```
go get github.com/Gonzih/replay-nginx-log
```

## Usage

```
Usage of replay-nginx-log:
  -debug
        Print extra debugging information
  -file string
        Log file name to read. Read from STDIN if file name is '-' (default "dummy")
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
replay-nginx-log --file my-acces.log --debug --log out.log
```

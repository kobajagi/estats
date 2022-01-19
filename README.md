## estats - A minimal metric retrieval system

A CLI tool which polls:
- Load average values

- Derived CPU percentage values
- Network interface statistics
- Disk partition usage in percent
- Memory usage in percent

Linux only.
May work on other *nix, but untested.

#### Build:

This is a Go tool with no external dependencies. To build the project run:

```
git checkout https://github.com/kobajagi/estats.git && cd estats
go build
```

#### How to use:

```
Usage:
  estats [OPTION...]
Options:
  -i int
    	interval in seconds at which to poll (default 5)
  -n string
    	network interface to poll
  -p string
    	disk partition to poll
```

If `-n` or `-p` options are not specified, corresponding metrics will not print. To specify disk partition, specify any file or folder path within that partition (as expected by `statfs`).

Tool output is in (compact) JSON format, piped through `jq` it looks like this (single poll):

```json
{
  "timestamp": 1642633633,
  "metrics": [
    {
      "name": "disk usage (%)",
      "metric": 8
    },
    {
      "name": "loadavg 1m",
      "metric": "0.29"
    },
    {
      "name": "loadavg 5m",
      "metric": "0.24"
    },
    {
      "name": "loadavg 15m",
      "metric": "0.33"
    },
    {
      "name": "cpu usage (%)",
      "metric": 5
    },
    {
      "name": "download (bytes)",
      "metric": 494800992
    },
    {
      "name": "upload (bytes)",
      "metric": 20615508
    },
    {
      "name": "memory usage (%)",
      "metric": 39
    }
  ]
}
```

Pressure
====

Simple multithreaded command line load testing tool written in Go

Please see https://github.com/rha7/pressure/releases for binaries.

```
Usage of pressure:
  -H value
    	header, can be specified multiple times (default apptypes.Headers{})
  -X string
    	requests' HTTP method to use (default "GET")
  -c uint
    	concurrent requests, minimum is 3 (default 10)
  -d string
    	data to be sent as body in request
  -l string
    	logging level (default "info")
  -n uint
    	total number of requests, mininum is number of concurrent threads (default 100)
  -p string
    	proxy to use, http assumed, scheme determines type, http or socks5
  -r	reuse connections (default false)
  -t uint
    	requests timeout (default 60)
```

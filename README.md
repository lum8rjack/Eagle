# Eagle

## Overview
Eagle is a quick and portable TCP scanner written in Go. It provides some basic functionality like reading a list of hosts, saving the results to a file, and adjusting the timeout and number of hosts to scan concurrently.

## Requirements
Eagle is written in Go and only uses the standard library. Once Go is installed, you can build the binary using Make.

Example:
```bash
make linux
```

## Examples
Eagle requires you to provide either a host or file containing hosts to scan, and which ports to scan for.

```bash
./eagle-Linux64.bin 
  -iL string
    	Input file containing a list of hosts to scan
  -ip string
    	Single host or CIDR address to scan
  -open
    	Only output open ports (default false)
  -output string
    	File to save the results
  -p string
    	Ports to scan (comma separated)
  -threads int
    	Number of hosts to scan concurrently (default 20)
  -timeout int
    	Timeout in seconds for each scanned port (default 5)

You must provide either a host to scan or list of hosts
```
Below is an example of scanning a /24 subnet and only showing hosts with port 22 open.
```bash
./eagle-Linux64.bin -p 22 -open -ip 10.10.1.0/24
Eagle started scan at 2022-03-26 09:43:05
10.10.1.5:22 Open
10.10.1.2:22 Open
10.10.1.1:22 Open
10.10.1.102:22 Open
10.10.1.100:22 Open
Completed scan in 39.93 seconds
```

The same scan with the number of concurrent IPs increased from 20 to 50 and timeout decreased from 5 to 3.
```bash
./eagle-Linux64.bin -p 22 -open -ip 10.10.1.0/24 -threads 50 -timeout 3
Eagle started scan at 2022-03-26 09:46:04
10.10.1.5:22 Open
10.10.1.2:22 Open
10.10.1.1:22 Open
10.10.1.102:22 Open
10.10.1.100:22 Open
Completed scan in 15.01 seconds
```

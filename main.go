package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"eagle/modules"
)

func main() {
	single_host := flag.String("ip", "", "Single host or CIDR address to scan")
	input_list := flag.String("iL", "", "Input file containing a list of hosts to scan")
	output_file := flag.String("output", "", "File to save the results")
	only_open := flag.Bool("open", false, "Only output open ports (default false)")
	selected_ports := flag.String("p", "", "Ports to scan (comma separated)")
	threads := flag.Int("threads", 25, "Number of hosts to scan concurrently")
	timeout := flag.Int("timeout", 3, "Timeout in seconds for each scanned port")

	flag.Parse()

	if *input_list == "" && *single_host == "" {
		flag.PrintDefaults()
		fmt.Println("\nYou must provide either a host to scan or list of hosts")
		os.Exit(1)
	}

	if *selected_ports == "" {
		flag.PrintDefaults()
		fmt.Println("\nYou must provide port(s) to scan")
		os.Exit(1)
	}

	// Setup the Logger
	logger := modules.Logger{
		SaveOutput: false,
		StartTime:  time.Now(),
		EndTime:    time.Now(),
	}

	if *output_file != "" {
		logger.OutputFilename = *output_file
		logger.SaveOutput = true
	}

	modules.SCANTIMEOUT = *timeout

	portsString := strings.Split(*selected_ports, ",")
	portsInt, err := modules.CheckPorts(portsString)
	if err != nil {
		fmt.Println("Error parsing list of ports")
		os.Exit(1)
	}
	modules.NUMIPS = len(portsInt)

	var ipList []string

	// Get IPs from file or single host
	if *input_list != "" {
		ipList = modules.ReadIPList(*input_list)
	} else {
		ipList = modules.CheckCIDR(*single_host)
	}

	// Setup number of go routines
	var wg sync.WaitGroup
	sem := make(chan int, *threads)

	// Keep track of open ports
	openPorts := 0

	// Start scanning
	logger.Start()

	// Loop through IPs to scan
	for _, i := range ipList {
		wg.Add(1)
		sem <- 1

		go func(i string) {
			defer wg.Done()
			host := modules.InitialScan(i, portsInt)

			// Print the output from the scans
			for _, p := range host.Ports {
				var data string
				if *only_open {
					if p.State == "Open" {
						data = fmt.Sprintf("%s:%d %s\n", host.Hostname, p.Port, p.State)
						openPorts++
					}
				} else {
					data = fmt.Sprintf("%s:%d %s\n", host.Hostname, p.Port, p.State)
				}
				logger.WriteToFile(data)
			}
			<-sem
		}(i)
	}

	// Wait for scans to be done
	wg.Wait()
	close(sem)

	// Scanning completed
	logger.Stop()

	// Log how many ports were open
	fmt.Printf("Number of open ports: %v\n", openPorts)

}

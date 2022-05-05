package modules

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var SCANTIMEOUT = 30

type Host struct {
	Hostname string
	Ports    []ScanResult
}

type ScanResult struct {
	Port  int
	State string
}

// Check that the CIDR/IP only contains
// specific characters/numbers, is less
// than 19 characters, and is a valid format
func isValidStringCIDR(cidr string) bool {
	// Must contain a /, and . and be less than 19 characters
	if !strings.Contains(cidr, "/") || !strings.Contains(cidr, ".") || len(cidr) > 18 {
		return false
	}

	// Keep track of the number of
	// . and / characters
	periods := 0
	slash := 0

	// Only valid characters
	valid := "0123456789./"
	s := strings.Split(cidr, "")
	if s[0] == "0" {
		return false
	}
	for _, c := range s {
		if !strings.Contains(valid, c) {
			return false
		}
		if c == "." {
			periods++
		} else if c == "/" {
			slash++
		}
	}

	// Make sure there are 3 periods and 1 slash
	if slash != 1 || periods != 3 {
		return false
	}

	return true
}

// Check if the provide string is a CIDR range
func CheckCIDR(line string) []string {
	var ips []string

	// Remove http/https if someone added to the front of the string
	line = strings.ToLower(line)
	if strings.HasPrefix(line, "https://") {
		line = strings.TrimPrefix(line, "https://")
	} else if strings.HasPrefix(line, "http://") {
		line = strings.TrimPrefix(line, "http://")
	}

	// Check if single host contains '/' and the length
	// is less than 18 (ex. 192.168.100.100/24)
	// Convert CIDR to list of IPs
	if isValidStringCIDR(line) {
		ips = cidrHosts(line)
	} else if line == "" {
		return ips
	} else if !strings.Contains(line, ".") {
		return ips
	} else {
		ips = append(ips, line)
	}

	return ips
}

// Take a filename and each line for an IP or CIDR range
func ReadIPList(file string) []string {
	var IPs []string

	// Open input file
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error reading input file")
		os.Exit(1)
	}
	defer f.Close()

	// Create a scanner to read the file
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	// Scan each line of the file and add to array
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	// Loop through each line of the file to check for CIDR
	for _, line := range fileLines {
		// Check if any of the lines are CIDR
		IPs = append(IPs, CheckCIDR(line)...)
	}

	return IPs
}

// FIXME: Need to make sure we don't add a port more than once
func CheckPorts(ports []string) ([]int, error) {
	var results []int

	if len(ports) == 0 {
		return results, errors.New("empty list of ports")
	}

	for _, p := range ports {
		portInt, err := strconv.Atoi(p)
		if err == nil {
			// Make sure it is a valid port number
			if portInt >= 1 && portInt <= 65535 {
				results = append(results, portInt)
			}
		}
	}

	return results, nil
}

// Convert a CIDR range to a list of IPs
func cidrHosts(netw string) []string {
	// convert string to IPNet struct
	_, ipv4Net, err := net.ParseCIDR(netw)
	if err != nil {
		return []string{}
	}
	// Convert IPNet struct mask and address to uint32
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	// Find the start IP address
	start := binary.BigEndian.Uint32(ipv4Net.IP)
	// Find the final IP address
	finish := (start & mask) | (mask ^ 0xffffffff)
	// Make a slice to return host addresses
	var hosts []string

	// Limit how many IPs we can scan
	// We can only take a /16
	n := (int)(finish - start)
	max := 65536
	if n > max {
		return []string{}
	}

	// Loop through addresses as uint32
	// Used "start + 1" and "finish - 1" to discard the network and broadcast addresses.
	for i := start + 1; i <= finish-1; i++ {
		// convert back to net.IPs
		// Create IP address of type net.IP. IPv4 is 4 bytes, IPv6 is 16 bytes.
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		hosts = append(hosts, ip.String())
	}
	// Return a slice of strings containing IP addresses
	fmt.Printf("Len: %d\n", len(hosts))
	return hosts
}

// Scan a one host for a single port and return Open/Closed
func scanPort(protocol, hostname string, port int) ScanResult {
	result := ScanResult{Port: port}

	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, time.Duration(SCANTIMEOUT)*time.Second)

	if err != nil {
		result.State = "Closed"
		return result
	}
	defer conn.Close()

	result.State = "Open"
	return result
}

// Take a single host and loop through the ports to scan
func InitialScan(hostname string, ports []int) Host {
	host := Host{Hostname: hostname}

	for _, p := range ports {
		host.Ports = append(host.Ports, scanPort("tcp", hostname, p))
	}

	return host
}

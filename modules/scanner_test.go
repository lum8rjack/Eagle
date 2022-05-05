package modules

import (
	"testing"
)

// TestCheckCIDR_Valid calls scanner.CheckCIDR
// with valid CIDR range
func TestCheckCIDR_Valid(t *testing.T) {

	c := "10.10.10.0/24"
	n := 254

	ips := CheckCIDR(c)
	if len(ips) != n {
		t.Fatalf("CheckCIDR(%s) response length = %d, should equal %d\n", c, len(ips), n)
	}
}

// TestCheckCIDR_Empty calls scanner.CheckCIDR
// with empty string
func TestCheckCIDR_Empty(t *testing.T) {

	c := ""
	n := 0

	ips := CheckCIDR(c)
	if len(ips) != n {
		t.Fatalf("CheckCIDR(%s) response length = %d, should equal %d\n", c, len(ips), n)
	}
}

// TestCheckCIDR_SingleIP calls scanner.CheckCIDR
// with a single IP instead of a CIDR range
func TestCheckCIDR_SingleIP(t *testing.T) {

	c := "10.10.10.10"
	n := 1

	ips := CheckCIDR(c)
	if len(ips) != n {
		t.Fatalf("CheckCIDR(%s) response length = %d, should equal %d\n", c, len(ips), n)
	}
}

// TestCheckCIDR_RandomString calls scanner.CheckCIDR
// with a random string
func TestCheckCIDR_RandonString(t *testing.T) {

	c := "string"
	n := 0

	ips := CheckCIDR(c)
	if len(ips) != n {
		t.Fatalf("CheckCIDR(%s) response length = %d, should equal %d\n", c, len(ips), n)
	}
}

// Fuzz the CheckCIDR function which takes
// one argument of an array of strings (ports)
// and converts it to an array of ints
func FuzzCheckCIDR(f *testing.F) {
	f.Add("10.10.10.0/24")
	f.Fuzz(func(t *testing.T, s string) {
		CheckCIDR(s)
	})
}

// TestCheckPorts_InvalidPort calls scanner.ChecPorts with
// a list of ports that contains a number
// outside of a valid range
func TestCheckPorts_InvalidPort(t *testing.T) {

	ports := []string{"22", "80", "-1", "659999"}
	n := 2

	p, _ := CheckPorts(ports)
	if len(p) != n {
		t.Fatalf("CheckPorts response length = %d, should equal %d\n", len(ports), n)
	}
}

// TestCheckPorts_Empty calls scanner.ChecPorts with
// an empty list of ports
func TestCheckPorts_Empty(t *testing.T) {

	ports := []string{}
	n := 0

	p, _ := CheckPorts(ports)
	if len(p) != n {
		t.Fatalf("CheckPorts response length = %d, should equal %d\n", len(ports), n)
	}
}

// TestCheckPorts_RandomString calls scanner.ChecPorts with
// a random string that is not a number
func TestCheckPorts_RandomString(t *testing.T) {

	ports := []string{"a", "https"}
	n := 0

	p, _ := CheckPorts(ports)
	if len(p) != n {
		t.Fatalf("CheckPorts response length = %d, should equal %d\n", len(ports), n)
	}
}

// Fuzz the CheckPorts function which takes
// one argument of an array of strings (ports)
// and converts it to an array of ints
func FuzzCheckPorts(f *testing.F) {
	f.Add("443")
	f.Fuzz(func(t *testing.T, s string) {
		p, err := CheckPorts([]string{s})
		if err != nil {
			t.Errorf("Fuzzing error: %s\n", err)
		} else if len(p) > 1 {
			t.Errorf("Fuzzing error wrong length (%d) with input []string{\"%s\"}\n", len(p), s)
		}
	})
}

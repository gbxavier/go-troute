package main

import (
	"fmt"
//	"errors"
	"net"
//	"syscall"
	"time"
	"flag"
)

// Hop Struct represents a single hop between a source and a destination.
type Hop struct {
	TTL       int
	AddrIP    net.IP
	AddrDNS   string
	Latency   time.Duration
	Err       error
}




func main() {

	// Default values, can be changed passing flags to the command.
	const DefaultMaxTTL   int = 30
	const DefaultFirstHop int = 0
	const DefatultProbes  int = 1


	// The Default values are setted above, but the user is able to change this values passing the respective flags to the command 
	// maxTTL is equals the last TTL used on calls
	// firstHop is the TTL used in the first call
	// probes is the number of calls executed for the same TTL
	var maxTTL = flag.Int("m", DefaultMaxTTL, "Set the max TTL (Time To Live) (default is 30)")
	var firstHop = flag.Int("f", DefaultFirstHop, "Set the first used Time-To-Live, e.g. the first hop (default is 1)")
	var probes = flag.Int("p", DefatultProbes, "Set the number of probes per 'TTL'(default is one probe).")

	flag.Parse()
	host := flag.Arg(0)
	
	ipAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return
	}

	fmt.Printf("Tracing route to %v [%v], over a maximum of %v hops, starting from %v:", host, ipAddr, *maxTTL, *firstHop)
	



}

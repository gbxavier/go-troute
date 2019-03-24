package main

import (
	"fmt"
	"flag"
	"errors"
	"math/rand"
	"net"
	"os"
	"time"
	icmp "golang.org/x/net/icmp"
	ipv4 "golang.org/x/net/ipv4"
	ipv6 "golang.org/x/net/ipv6"
)

// Hop Struct represents a single hop between a source and a destination.
type Hop struct {
	TTL     int
	AddrIP  net.IP
	AddrDNS []string
	Latency time.Duration
	Err     error
}

func createICMPEcho(ICMPTypeEcho icmp.Type) (req []byte, err error) {
	echo := icmp.Message{
		Type: ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Int(),
			Seq:  1,
			Data: []byte("TABS"),
		}}

	req, err = echo.Marshal(nil)
	return
}

func callHop(host string, ttl int, req []byte, proto string, dialProto string, dialDest string, ipVersion int, listenAddress string, timeout int) (currentHop Hop, err error){
	

	// Opening connection to host
	conn, err := net.Dial(dialProto, dialDest)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Opening outbound ipv4 or ipv6 connection, for ipv4 or ipv6 protocols, respectively
	if ipVersion == 4 {

		newConn := ipv4.NewConn(conn)
		if err = newConn.SetTTL(ttl); err != nil {
			fmt.Println(err)
			return
		}

	} else {

		newConn := ipv6.NewConn(conn)
		if err = newConn.SetHopLimit(ttl); err != nil {
			fmt.Println(err)
			return
		}

	}

	// Opening Inbound ICMP Listener
	packetConn, err := icmp.ListenPacket("ip"+fmt.Sprintf("%d", ipVersion)+":"+"icmp", listenAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer packetConn.Close()

	// Starting counter and sending request
	start := time.Now()
	_, err = conn.Write(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = packetConn.SetDeadline(start.Add(time.Second * time.Duration(timeout))); err != nil {
		fmt.Println(err)
		return
	}

	// Reading ICMP packet, if exists.
	readBytes := make([]byte, 1500)                     // 1500 Bytes ethernet MTU
	_, sAddr, connErr := packetConn.ReadFrom(readBytes) 

	latency := time.Since(start)
	
	currentHop = Hop{
		TTL     : ttl,
		Latency : latency,
		Err     : connErr,
	}	
	
	if connErr == nil {
		currentHop.AddrIP = net.ParseIP(sAddr.String())
		if currentHop.AddrIP == nil {
			currentHop.Err = errors.New("Timeout")
		}else{
			currentHop.AddrDNS, _ = net.LookupAddr(currentHop.AddrIP.String())
		}
	}



	return currentHop, nil
}

func printHop(hop Hop){

	if hop.AddrIP == nil{
		fmt.Printf("%d - * - Request timed out\n", hop.TTL)
	}else{
		fmt.Printf("%d - %v %v - time elapsed: %v\n", hop.TTL, hop.AddrIP, hop.AddrDNS, hop.Latency)
	}
		

}

// Hop set with all the hops, returned from this.
func traceRoute(host string, maxTTL *int, firstHop *int, proto string, ipVersion int) (hopSet []Hop) {

	var dialProto string
	var dialDest = host
	var listenAddress string
	var req = []byte{}
	const DefaultTimeoutS int = 3
	ttl := *firstHop
	var found = false
	

	// Try to resolve the host provided, if name returns the ip address
	ipAddr, err := net.ResolveIPAddr(fmt.Sprintf("ip%d", ipVersion), host)
	if err != nil {
		fmt.Println("Error resolving IP")
		os.Exit(1)
		return
	}

	// User feedback of what will happen
	fmt.Printf("Tracing route to %v [%v], over a maximum of %d hops, starting from %d:\n\n", host, ipAddr, *maxTTL, *firstHop)

	// Configuring for UDP packets or ICMP packets
	if proto == "udp" {

		// Sending UDP packets
		dialProto = "udp" + fmt.Sprintf("%d", ipVersion)
		dialDest += ":33454"
		req = []byte("TABS")

	} else {

		// Sending ICMP packets
		dialProto = "ip" + fmt.Sprintf("%d", ipVersion) + ":" + proto // icmp

		// Creating the request for the current IP Version
		if ipVersion == 4 {
			req, err = createICMPEcho(ipv4.ICMPTypeEcho)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}
		} else {
			req, err = createICMPEcho(ipv6.ICMPTypeEchoRequest)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}
		}
	}

	if ipVersion == 4 {
		listenAddress = "0.0.0.0"
	} else {
		listenAddress = "::0"
	}

	for i := ttl; i <= *maxTTL; i++ {

		current, err := callHop(host, i, req, proto, dialProto, dialDest, ipVersion, listenAddress, DefaultTimeoutS)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		printHop(current)
	
		hopSet = append(hopSet, current)

		if current.AddrIP.String() == ipAddr.IP.String(){
			found = true
			break
		}

	}

	if found == false{
		fmt.Println("Not Found, please consider increase TTL")
	}else{
		fmt.Println("\n Trace Complete")
	}

	return hopSet
}

func main() {

	// Default values, can be changed passing flags to the command.
	const DefaultMaxTTL int = 30
	const DefaultFirstHop int = 1

	// The Default values are setted above, but the user is able to change this values passing the respective flags to the command
	// maxTTL is equals the last TTL used on calls
	// firstHop is the TTL used in the first call
	var maxTTL = flag.Int("m", DefaultMaxTTL, "Set the max TTL (Time To Live) (default is 30)")
	var firstHop = flag.Int("f", DefaultFirstHop, "Set the first used Time-To-Live, e.g. the first hop (default is 1)")
	var ipv6 = flag.Bool("ipv6", false, "Set to IPV6.")
	var udp = flag.Bool("udp", false, "Set to UDP Mode.")
	flag.Parse()

	host := flag.Arg(0)

	var ipVersion = 4
	if *ipv6 == true {
		fmt.Println("IPV6 Mode")
		ipVersion = 6
	}

	var proto = "icmp"
	if *udp == true {
		fmt.Println("UDP Mode")
		proto = "udp"
	}

	// Check if the target host is provided, and exits if not
	if len(host) == 0 {
		fmt.Println("Please, specify a host")
		os.Exit(1)
		return
	}

	// The traceroute execution itself
	hopSet := traceRoute(host, maxTTL, firstHop,  proto, ipVersion)

	
	//Removing timed out hops to calculate response beetween consecutive hops
	for i, hop := range hopSet {
        if hop.AddrIP == nil {
            hopSet = append(hopSet[:i], hopSet[i+1:]...)
        }
    }

	var maxLeft Hop
	var maxRight Hop
	var maxLatency time.Duration

	// This section compares the time between hops, and store these.
	for indexLeft, indexRight := 0, 1; indexRight < len(hopSet) ; indexLeft, indexRight = indexLeft+1, indexRight+1{
		
		if (hopSet[indexRight].Latency - hopSet[indexLeft].Latency) > maxLatency {
			
			maxLatency = hopSet[indexRight].Latency - hopSet[indexLeft].Latency
			maxLeft = hopSet[indexLeft]
			maxRight = hopSet[indexRight]

		}
		
	}

	fmt.Printf("\nThe largest difference in response time between consecutive hops is %v\n", maxLatency)
	fmt.Println("Between hops: ")
	printHop(maxLeft)
	printHop(maxRight)

}

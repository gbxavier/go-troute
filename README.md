# Traceroute

In [computing](https://en.wikipedia.org/wiki/Computing "Computing"), **`traceroute`** and **`tracert`** are [computer network](https://en.wikipedia.org/wiki/Computer_network "Computer network") diagnostic [commands](https://en.wikipedia.org/wiki/Command_(computing) "Command (computing)") for displaying the route (path) and measuring transit delays of [packets](https://en.wikipedia.org/wiki/Network_packet "Network packet") across an [Internet Protocol](https://en.wikipedia.org/wiki/Internet_Protocol "Internet Protocol") (IP) network.

This implementation of traceroute was written using [Go Lang](https://golang.org/), and the built-in libraries: *fmt*, *flag*, *errors*, *math/rand*, *net*, *os*, *time*, and the external libraries: golang.org/x/net/icmp, golang.org/x/net/ipv4, golang.org/x/net/ipv6. For testing, is used "*testing*".

 - **Jenkins** as CI/CD Tool. Builds, Tests, and Publish.

# Usage
To use this tool, the first step is clone this repo. Inside the repo execute the commands below:

```sh
go get -d -v ./...
go build -o traceroute-gbxavier11
./traceroute-gbxavier11 [options] $URL 
```
Option list:

 - -f int
	 - Set the first used Time-To-Live, e.g. the first hop (default 1)
 - -ipv6
	 - Set to IPV6 Mode.
 - -m int
	 - Set the max TTL (Time To Live) (default 30)
 - -udp
	 -  Set to UDP Mode (default ICMP).

Example:

```sh
./traceroute-gbxavier11 -m 50 -udp www.google.com 
```
Note[1]: Depending on your operating system, you need to choice between UDP or ICMP to run this command.
Note[2]: Make sure you have the right permissions to run the command.
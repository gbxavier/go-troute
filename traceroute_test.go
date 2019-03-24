package main 

import (
	"testing"
)

func TestEmpty(t *testing.T) {
	
	maxTTL := 10
	firstHop := 1
	proto := "icmp"
	ipVersion := 4
	got := traceRoute("", &maxTTL, &firstHop, proto, ipVersion)
	
	if len(got) != 0 {
		t.Errorf("traceRoute() = %v; want []", got)
	}
	
}
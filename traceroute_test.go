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

func TestGoogle(t *testing.T){
	
	maxTTL := 10
	firstHop := 1
	got := traceRoute("www.google.com", &maxTTL, &firstHop, "icmp", 4)
	
	if len(got) == 0 {
		t.Errorf("len(traceRoute()) == %v; want x > 0", len(got))
	}

}

func TestGoogleLowTTL(t *testing.T){
	
	maxTTL := 2
	firstHop := 1
	got := traceRoute("www.google.com", &maxTTL, &firstHop, "icmp", 4)
	
	if len(got) != 2 {
		t.Errorf("len(traceRoute()) == %v; want 2", len(got))
	}
	
}
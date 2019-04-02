package main

import (
	"fmt"
	"net"
)

const PROBE_DNS_NAME = "edgy"

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	addrs, err := net.LookupHost(PROBE_DNS_NAME)
	panicOnErr(err)
	for _, addr := range addrs {
		metrics, err := fetchMetrics("http://" + addr + ":8123/")
		panicOnErr(err)
		fmt.Println(metrics)
	}
}

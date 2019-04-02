package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	probeDNS := os.Getenv("PROBE_DNS_NAME")
	if len(probeDNS) == 0 {
		panic("PROBE_DNS_NAME is not set")
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		addrs, err := net.LookupHost(probeDNS)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ret := ""
		for _, addr := range addrs {
			metrics, err := fetchMetrics("http://" + addr + ":8123/")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			ret += toPromFormat(metrics)
		}
		_, err = fmt.Fprintln(w, ret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	if err := http.ListenAndServe(":8321", nil); err != nil {
		panicOnErr(err)
	}
}

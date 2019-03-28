package main

import (
	"fmt"
	"net/http"

	"github.com/shirou/gopsutil/mem"
)

func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		v, _ := mem.VirtualMemory()
		fmt.Fprintln(w, v)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

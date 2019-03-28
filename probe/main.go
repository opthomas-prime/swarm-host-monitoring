package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	docker "github.com/docker/docker/client"

	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

type NodeInfo struct {
	Name string `json:"name"`
}

func main() {
	cli, err := docker.NewEnvClient()
	panicOnErr(err)

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		dockerInfo, err := cli.Info(context.Background())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		info := NodeInfo{dockerInfo.Name}
		bytes, err := json.Marshal(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = fmt.Fprintln(w, string(bytes))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/mem", func(w http.ResponseWriter, r *http.Request) {
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = fmt.Fprintln(w, memInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/load", func(w http.ResponseWriter, r *http.Request) {
		loadInfo, err := load.Avg()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = fmt.Fprintln(w, loadInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":8123", nil); err != nil {
		panicOnErr(err)
	}
}

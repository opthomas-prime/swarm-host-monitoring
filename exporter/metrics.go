package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type nodeInfo struct {
	Name string `json:"name"`
}

type nodeLoad struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type nodeMem struct {
	Total     int64 `json:"total"`
	Available int64 `json:"available"`
	Used      int64 `json:"used"`
}

type nodeMetrics struct {
	info nodeInfo
	load nodeLoad
	mem  nodeMem
}

func fetchInfo(addr string) (nodeInfo, error) {
	resp, err := http.Get(addr + "info")
	info := nodeInfo{}
	if err != nil {
		return info, err
	}
	if resp.StatusCode != 200 {
		return info, errors.New("unexpected HTTP status")
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&info)
	return info, err
}

func fetchLoad(addr string) (nodeLoad, error) {
	resp, err := http.Get(addr + "load")
	load := nodeLoad{}
	if err != nil {
		return load, err
	}
	if resp.StatusCode != 200 {
		return load, errors.New("unexpected HTTP status")
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&load)
	return load, err
}

func fetchMem(addr string) (nodeMem, error) {
	resp, err := http.Get(addr + "mem")
	mem := nodeMem{}
	if err != nil {
		return mem, err
	}
	if resp.StatusCode != 200 {
		return mem, errors.New("unexpected HTTP status")
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&mem)
	return mem, err
}

func fetchMetrics(addr string) (nodeMetrics, error) {
	metrics := nodeMetrics{}
	info, err := fetchInfo(addr)
	if err != nil {
		return metrics, err
	}
	metrics.info = info

	load, err := fetchLoad(addr)
	if err != nil {
		return metrics, err
	}
	metrics.load = load

	mem, err := fetchMem(addr)
	if err != nil {
		return metrics, err
	}
	metrics.mem = mem

	return metrics, nil
}

func toPromFormat(metrics nodeMetrics) string {
	prefix := "swarm_host_monitoring_"
	instance := metrics.info.Name
	value := ""
	value += prefix + "load1{instance=\"" + instance + "\"} " + fmt.Sprintf("%f", metrics.load.Load1) + "\n"
	value += prefix + "load5{instance=\"" + instance + "\"} " + fmt.Sprintf("%f", metrics.load.Load5) + "\n"
	value += prefix + "load15{instance=\"" + instance + "\"} " + fmt.Sprintf("%f", metrics.load.Load15) + "\n"
	value += prefix + "mem_total{instance=\"" + instance + "\"} " + strconv.FormatInt(metrics.mem.Total, 10) + "\n"
	value += prefix + "mem_available{instance=\"" + instance + "\"} " + strconv.FormatInt(metrics.mem.Available, 10) + "\n"
	value += prefix + "mem_used{instance=\"" + instance + "\"} " + strconv.FormatInt(metrics.mem.Used, 10) + "\n"
	return value
}

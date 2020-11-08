package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

type AvgStat struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

func (l AvgStat) String() string {
	s, _ := json.Marshal(l)
	return string(s)
}

type MiscStat struct {
	ProcsTotal   int `json:"procsTotal"`
	ProcsRunning int `json:"procsRunning"`
	ProcsBlocked int `json:"procsBlocked"`
	Ctxt         int `json:"ctxt"`
}

func (m MiscStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

func (ps *PSUtils) ArgLoad() (*AvgStat, error) {
	values, err := ps.readLoadAvgFromFile()
	if err != nil {
		return nil, err
	}

	load1, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return nil, err
	}
	load5, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return nil, err
	}
	load15, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		return nil, err
	}

	ret := &AvgStat{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}

	return ret, nil
}

func (ps *PSUtils) MiscLoad() (*MiscStat, error) {

}

func (ps *PSUtils)readLoadAvgFromFile() ([]string, error) {
	loadavgFilename := "/proc/loadavg"
	line, err := ps.FileContent(loadavgFilename)
	if err != nil {
		return nil, err
	}

	values := strings.Fields(line)
	return values, nil
}



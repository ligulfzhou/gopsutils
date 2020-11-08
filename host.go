package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type InfoStat struct {
	Hostname             string `json:"hostname"`
	Uptime               uint64 `json:"uptime"`
	BootTime             uint64 `json:"bootTime"`
	Procs                uint64 `json:"procs"`           // number of processes
	OS                   string `json:"os"`              // ex: freebsd, linux
	Platform             string `json:"platform"`        // ex: ubuntu, linuxmint
	PlatformFamily       string `json:"platformFamily"`  // ex: debian, rhel
	PlatformVersion      string `json:"platformVersion"` // version of the complete OS
	KernelVersion        string `json:"kernelVersion"`   // version of the OS kernel (if available)
	KernelArch           string `json:"kernelArch"`      // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
	VirtualizationSystem string `json:"virtualizationSystem"`
	VirtualizationRole   string `json:"virtualizationRole"` // guest or host
	HostID               string `json:"hostId"`             // ex: uuid
}

type lsbStruct struct {
	ID          string
	Release     string
	Codename    string
	Description string
}

func (lsb *lsbStruct) String() string {
	s, _ := json.Marshal(lsb)
	return string(s)
}

type UserStat struct {
	User     string `json:"user"`
	Terminal string `json:"terminal"`
	Host     string `json:"host"`
	Started  int    `json:"started"`
}

type TemperatureStat struct {
	SensorKey   string  `json:"sensorKey"`
	Temperature float64 `json:"temperature"`
	High        float64 `json:"sensorHigh"`
	Critical    float64 `json:"sensorCritical"`
}

func (h InfoStat) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (u UserStat) String() string {
	s, _ := json.Marshal(u)
	return string(s)
}

func (t TemperatureStat) String() string {
	s, _ := json.Marshal(t)
	return string(s)
}

func (ps PSUtils) getlsbStruct() (*lsbStruct, error) {
	// var err error
	ret := &lsbStruct{}
	if ps.FileExists("/etc/lsb-release") {
		str, err := ps.Exec("cat /etc/lsb-release")
		if err != nil {
			return nil, err
		}

		seq := "\n"
		if strings.Contains(str, "\r\n") {
			seq = "\r\n"
		}
		strs := strings.Split(str, seq)
		for _, s := range strs {
			field := strings.Split(s, "=")
			if len(field) < 2 {
				continue
			}
			switch field[0] {
			case "DISTRIB_ID":
				ret.ID = StripString(field[1])
			case "DISTRIB_RELEASE":
				ret.Release = StripString(field[1])
			case "DISTRIB_CODENAME":
				ret.Codename = StripString(field[1])
			case "DISTRIB_DESCRIPTION":
				ret.Description = StripString(field[1])
			}
		}
	} else if ps.FileExists("/usr/bin/lsb_release") {
		str, err := ps.Exec("lsb_release")
		if err != nil {
			return nil, err
		}

		seq := "\n"
		if strings.Contains(str, "\r\n") {
			seq = "\r\n"
		}
		strs := strings.Split(str, seq)
		for _, s := range strs {
			field := strings.Split(s, "=")
			if len(field) < 2 {
				continue
			}
			switch field[0] {
			case "Distributor ID":
				ret.ID = StripString(field[1])
			case "Release":
				ret.Release = StripString(field[1])
			case "Codename":
				ret.Codename = StripString(field[1])
			case "Description":
				ret.Description = StripString(field[1])
			}
		}
	}

	return ret, nil
}

func (ps PSUtils) PlatformInformationWithContext(ctx context.Context) (platform string, family string, version string, err error) {
	// var err error
	lsb, err := ps.getlsbStruct()
	fmt.Println(lsb)

	return "", "", "", nil
}

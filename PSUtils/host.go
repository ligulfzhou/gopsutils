package PSUtils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	SysProductUUID            = "/sys/class/dmi/id/product_uuid"
	MachineID                 = "/etc/machine-id"
	ProcSysKernelRandomBootID = "/proc/sys/kernel/random/boot_id"
)

type HostInfoStat struct {
	Hostname             string `json:"hostname"`
	Uptime               int64  `json:"uptime"`
	BootTime             int64  `json:"bootTime"`
	Procs                int64  `json:"procs"`           // number of processes
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

func (h HostInfoStat) String() string {
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

func (ps *PSUtils) GetHostInfoStat() (*HostInfoStat, error) {
	var (
		err error
	)

	ret := &HostInfoStat{}
	ret.Platform, ret.PlatformFamily, ret.PlatformVersion, err = ps.PlatformInformation()
	if err != nil {
		fmt.Printf("get platform... failed: %s", err.Error())
	}

	ret.KernelVersion = ps.GetKernelVersion()
	ret.KernelArch = ps.GetKernalArch()

	return ret, nil
}

func (ps *PSUtils) getlsbStruct() (*lsbStruct, error) {
	// var err error
	ret := &lsbStruct{}
	if ps.FileExists("/etc/lsb-release") {
		str, err := ps.Exec("cat /etc/lsb-release")
		if err != nil {
			return nil, err
		}

		strs := SplitStringToLines(str)
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

		strs := SplitStringToLines(str)
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

func (ps *PSUtils) PlatformInformation() (platform string, family string, version string, err error) {
	lsb, err := ps.getlsbStruct()
	if err != nil {
		lsb = &lsbStruct{}
	}

	if ps.FileExists("/etc/oracle-release") {
		platform = "oracle"
		contents, err := ps.ReadLines("/etc/oracle-release")
		if err == nil {
			version = getRedhatishVersion(contents)
		}
	} else if ps.FileExists("/etc/enterprise-release") {
		platform = "oracle"
		contents, err := ps.ReadLines("/etc/enterprise-release")
		if err == nil {
			version = getRedhatishVersion(contents)
		}
	} else if ps.FileExists("/etc/slackware-version") {
		platform = "slackware"
		contents, err := ps.ReadLines("/etc/slackware-version")
		if err == nil {
			version = getSlackwareVersion(contents)
		}
	} else if ps.FileExists("/etc/debian_version") {
		if lsb.ID == "Ubuntu" {
			platform = "ubuntu"
			version = lsb.Release
		} else if lsb.ID == "LinuxMint" {
			platform = "linuxmint"
			version = lsb.Release
		} else {
			if ps.FileExists("/usr/bin/raspi-config") {
				platform = "raspbian"
			} else {
				platform = "debian"
			}
			contents, err := ps.ReadLines("/etc/debian_version")
			if err == nil && len(contents) > 0 && contents[0] != "" {
				version = contents[0]
			}
		}
	} else if ps.FileExists("/etc/redhat-release") {
		contents, err := ps.ReadLines("/etc/redhat-release")
		if err == nil {
			version = getRedhatishVersion(contents)
			platform = getRedhatishPlatform(contents)
		}
	} else if ps.FileExists("/etc/system-release") {
		contents, err := ps.ReadLines("/etc/system-release")
		if err == nil {
			version = getRedhatishVersion(contents)
			platform = getRedhatishPlatform(contents)
		}
	} else if ps.FileExists("/etc/gentoo-release") {
		platform = "gentoo"
		contents, err := ps.ReadLines("/etc/gentoo-release")
		if err == nil {
			version = getRedhatishVersion(contents)
		}
	} else if ps.FileExists("/etc/SuSE-release") {
		contents, err := ps.ReadLines("/etc/SuSE-release")
		if err == nil {
			version = getSuseVersion(contents)
			platform = getSusePlatform(contents)
		}
		// TODO: slackware detecion
	} else if ps.FileExists("/etc/arch-release") {
		platform = "arch"
		version = lsb.Release
	} else if ps.FileExists("/etc/alpine-release") {
		platform = "alpine"
		contents, err := ps.ReadLines("/etc/alpine-release")
		if err == nil && len(contents) > 0 && contents[0] != "" {
			version = contents[0]
		}
	} else if ps.FileExists("/etc/os-release") {
		pv := ps.GetOSRelease()
		if err == nil {
			platform = pv[0]
			version = pv[1]
		}
	} else if lsb.ID == "RedHat" {
		platform = "redhat"
		version = lsb.Release
	} else if lsb.ID == "Amazon" {
		platform = "amazon"
		version = lsb.Release
	} else if lsb.ID == "ScientificSL" {
		platform = "scientific"
		version = lsb.Release
	} else if lsb.ID == "XenServer" {
		platform = "xenserver"
		version = lsb.Release
	} else if lsb.ID != "" {
		platform = strings.ToLower(lsb.ID)
		version = lsb.Release
	}

	switch platform {
	case "debian", "ubuntu", "linuxmint", "raspbian":
		family = "debian"
	case "fedora":
		family = "fedora"
	case "oracle", "centos", "redhat", "scientific", "enterpriseenterprise", "amazon", "xenserver", "cloudlinux", "ibm_powerkvm":
		family = "rhel"
	case "suse", "opensuse", "sles":
		family = "suse"
	case "gentoo":
		family = "gentoo"
	case "slackware":
		family = "slackware"
	case "arch":
		family = "arch"
	case "exherbo":
		family = "exherbo"
	case "alpine":
		family = "alpine"
	case "coreos":
		family = "coreos"
	case "solus":
		family = "solus"
	}

	fmt.Println(platform, family, version, err)
	return platform, family, version, nil
}

func getSlackwareVersion(contents []string) string {
	c := strings.ToLower(strings.Join(contents, ""))
	c = strings.Replace(c, "slackware ", "", 1)
	return c
}

func getRedhatishVersion(contents []string) string {
	c := strings.ToLower(strings.Join(contents, ""))

	if strings.Contains(c, "rawhide") {
		return "rawhide"
	}
	if matches := regexp.MustCompile(`release (\d[\d.]*)`).FindStringSubmatch(c); matches != nil {
		return matches[1]
	}
	return ""
}

func getRedhatishPlatform(contents []string) string {
	c := strings.ToLower(strings.Join(contents, ""))

	if strings.Contains(c, "red hat") {
		return "redhat"
	}
	f := strings.Split(c, " ")

	return f[0]
}

func getSuseVersion(contents []string) string {
	version := ""
	for _, line := range contents {
		if matches := regexp.MustCompile(`VERSION = ([\d.]+)`).FindStringSubmatch(line); matches != nil {
			version = matches[1]
		} else if matches := regexp.MustCompile(`PATCHLEVEL = ([\d]+)`).FindStringSubmatch(line); matches != nil {
			version = version + "." + matches[1]
		}
	}
	return version
}

func getSusePlatform(contents []string) string {
	c := strings.ToLower(strings.Join(contents, ""))
	if strings.Contains(c, "opensuse") {
		return "opensuse"
	}
	return "suse"
}

func (ps *PSUtils) GetOSRelease() []string {
	var platform, version string
	contents, err := ps.ReadLines("/etc/os-release")
	if err != nil {
		return []string{"", ""}
	}

	for _, line := range contents {
		field := strings.Split(line, "=")
		if len(field) < 2 {
			continue
		}
		switch field[0] {
		case "ID": // use ID for lowercase
			platform = TrimQuotes(field[1])
		case "VERSION":
			version = TrimQuotes(field[1])
		}
	}
	return []string{platform, version}
}

func (ps *PSUtils) GetHostID() string {
	if ps.HostId != "" {
		return ps.HostId
	}

	hostId := ps.HostID()
	ps.HostId = hostId
	return hostId
}

func (ps *PSUtils) HostID() string {
	// In order to read this file, needs to be supported by kernel/arch and run as root
	// so having fallback is important
	if ps.FileExists(SysProductUUID) {
		s, err := ps.Exec(fmt.Sprintf("cat %s", SysProductUUID))
		if err == nil {
			return strings.ToLower(StripString(s))
		}
	}

	// Fallback on GNU Linux systems with systemd, readable by everyone
	if ps.FileExists(MachineID) {
		s, err := ps.Exec(fmt.Sprintf("cat %s", MachineID))
		if err == nil {
			hostId := StripString(s)
			if len(hostId) == 32 {
				return fmt.Sprintf("%s-%s-%s-%s-%s", hostId[0:8], hostId[8:12], hostId[12:16], hostId[16:20], hostId[20:32])
			}
		}
	}

	// Not stable between reboot, but better than nothing
	s, err := ps.Exec(fmt.Sprintf("cat %s", ProcSysKernelRandomBootID))
	if err == nil {
		return strings.ToLower(StripString(s))
	}

	return ""
}

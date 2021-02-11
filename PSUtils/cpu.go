package PSUtils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type CPUTimesStat struct {
	CPU       string  `json:"cpu"`
	User      float64 `json:"user"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	Nice      float64 `json:"nice"`
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guestNice"`
}

type CPUInfoStat struct {
	CPU        int32    `json:"cpu"`
	VendorID   string   `json:"vendorId"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int32    `json:"stepping"`
	PhysicalID string   `json:"physicalId"`
	CoreID     string   `json:"coreId"`
	Cores      int32    `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int32    `json:"cacheSize"`
	Flags      []string `json:"flags"`
	Microcode  string   `json:"microcode"`
}

type lastPercent struct {
	sync.Mutex
	lastCPUTimes    []CPUTimesStat
	lastPerCPUTimes []CPUTimesStat
}

func (ps *PSUtils) GetCPUStats() {

}

// number of Cores
// Counts returns the number of physical or logical cores in the system
func (ps *PSUtils) CPUCount(logical bool) (int, error) {
	if logical {
		ret := 0
		// https://github.com/giampaolo/psutil/blob/d01a9eaa35a8aadf6c519839e987a49d8be2d891/psutil/_pslinux.py#L599
		procCpuinfo := "/proc/cpuinfo"
		lines, err := ps.ReadLines(procCpuinfo)
		if err == nil {
			for _, line := range lines {
				line = strings.ToLower(line)
				if strings.HasPrefix(line, "processor") {
					ret++
				}
			}
		}
		if ret == 0 {
			procStat := "/proc/stat"
			lines, err = ps.ReadLines(procStat)
			if err != nil {
				return 0, err
			}
			for _, line := range lines {
				if len(line) >= 4 && strings.HasPrefix(line, "cpu") && '0' <= line[3] && line[3] <= '9' { // `^cpu\d` regexp matching
					ret++
				}
			}
		}
		return ret, nil
	}

	// physical cores
	// https://github.com/giampaolo/psutil/blob/122174a10b75c9beebe15f6c07dcf3afbe3b120d/psutil/_pslinux.py#L621-L629
	var threadSiblingsLists = make(map[string]bool)
	if files, err := ps.Glob("/sys/devices/system/cpu/cpu[0-9]*/topology/thread_siblings_list"); err == nil {
		for _, file := range files {
			lines, err := ps.ReadLines(file)
			if err != nil || len(lines) != 1 {
				continue
			}
			threadSiblingsLists[lines[0]] = true
		}
		ret := len(threadSiblingsLists)
		if ret != 0 {
			return ret, nil
		}
	}
	// https://github.com/giampaolo/psutil/blob/122174a10b75c9beebe15f6c07dcf3afbe3b120d/psutil/_pslinux.py#L631-L652
	filename := "/proc/cpuinfo"
	lines, err := ps.ReadLines(filename)
	if err != nil {
		return 0, err
	}
	mapping := make(map[int]int)
	currentInfo := make(map[string]int)
	for _, line := range lines {
		line = strings.ToLower(strings.TrimSpace(line))
		if line == "" {
			// new section
			id, okID := currentInfo["physical id"]
			cores, okCores := currentInfo["cpu cores"]
			if okID && okCores {
				mapping[id] = cores
			}
			currentInfo = make(map[string]int)
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		fields[0] = strings.TrimSpace(fields[0])
		if fields[0] == "physical id" || fields[0] == "cpu cores" {
			val, err := strconv.Atoi(strings.TrimSpace(fields[1]))
			if err != nil {
				continue
			}
			currentInfo[fields[0]] = val
		}
	}
	ret := 0
	for _, v := range mapping {
		ret += v
	}
	return ret, nil
}

func (ps *PSUtils) CPUInfo() ([]CPUInfoStat, error) {
	filename := "/proc/cpuinfo"
	lines, _ := ps.ReadLines(filename)

	var ret []CPUInfoStat
	var processorName string

	c := CPUInfoStat{CPU: -1, Cores: 1}
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])

		switch key {
		case "Processor":
			processorName = value
		case "processor":
			if c.CPU >= 0 {
				err := ps.finishCPUInfo(&c)
				if err != nil {
					return ret, err
				}
				ret = append(ret, c)
			}
			c = CPUInfoStat{Cores: 1, ModelName: processorName}
			t, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return ret, err
			}
			c.CPU = int32(t)
		case "vendorId", "vendor_id":
			c.VendorID = value
		case "cpu family":
			c.Family = value
		case "model":
			c.Model = value
		case "model name", "cpu":
			c.ModelName = value
			if strings.Contains(value, "POWER8") ||
				strings.Contains(value, "POWER7") {
				c.Model = strings.Split(value, " ")[0]
				c.Family = "POWER"
				c.VendorID = "IBM"
			}
		case "stepping", "revision":
			val := value

			if key == "revision" {
				val = strings.Split(value, ".")[0]
			}

			t, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return ret, err
			}
			c.Stepping = int32(t)
		case "cpu MHz", "clock":
			// treat this as the fallback value, thus we ignore error
			if t, err := strconv.ParseFloat(strings.Replace(value, "MHz", "", 1), 64); err == nil {
				c.Mhz = t
			}
		case "cache size":
			t, err := strconv.ParseInt(strings.Replace(value, " KB", "", 1), 10, 64)
			if err != nil {
				return ret, err
			}
			c.CacheSize = int32(t)
		case "physical id":
			c.PhysicalID = value
		case "core id":
			c.CoreID = value
		case "flags", "Features":
			c.Flags = strings.FieldsFunc(value, func(r rune) bool {
				return r == ',' || r == ' '
			})
		case "microcode":
			c.Microcode = value
		}
	}
	if c.CPU >= 0 {
		err := ps.finishCPUInfo(&c)
		if err != nil {
			return ret, err
		}
		ret = append(ret, c)
	}
	return ret, nil
}

func (ps *PSUtils) finishCPUInfo(c *CPUInfoStat) error {
	var lines []string
	var err error
	var value float64

	if len(c.CoreID) == 0 {
		lines, err = ps.ReadLines(sysCPUPath(c.CPU, "topology/core_id"))
		if err == nil {
			c.CoreID = lines[0]
		}
	}

	// override the value of c.Mhz with cpufreq/cpuinfo_max_freq regardless
	// of the value from /proc/cpuinfo because we want to report the maximum
	// clock-speed of the CPU for c.Mhz, matching the behaviour of Windows
	lines, err = ps.ReadLines(sysCPUPath(c.CPU, "cpufreq/cpuinfo_max_freq"))
	// lines, err = common.ReadLines(sysCPUPath(c.CPU, "cpufreq/cpuinfo_max_freq"))
	// if we encounter errors below such as there are no cpuinfo_max_freq file,
	// we just ignore. so let Mhz is 0.
	if err != nil || len(lines) == 0 {
		return nil
	}
	value, err = strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return nil
	}
	c.Mhz = value / 1000.0 // value is in kHz
	if c.Mhz > 9999 {
		c.Mhz = c.Mhz / 1000.0 // value in Hz
	}
	return nil
}

func sysCPUPath(cpu int32, relPath string) string {
	return fmt.Sprintf("/sys/devices/system/cpu/cpu%d/%s", cpu, relPath)
}

func (ps *PSUtils) CpuTimes() {

}

/*
     user    nice   system  idle      iowait irq   softirq  steal  guest  guest_nice
cpu  74608   2520   24433   1117073   6176   4054  0        0      0      0

Time units are in USER_HZ(1/100 second)

(user, nice, system, idle, iowait, irq, softirq [steal, [guest, [guest_nice]]])
Last 3 fields may not be available on all Linux kernel versions.
*/
type CpuTimes struct {
	User      int64
	Nice      int64
	System    int64
	Idle      int64
	Iowait    int64
	Irq       int64
	Softirq   int64
	Steal     int64
	Guest     int64
	GuestNice int64
}

func (ct CpuTimes) String() string {
	s, _ := json.Marshal(ct)
	return string(s)
}

func (ps *PSUtils) TotalCpuTimes() CpuTimes {
	var ret CpuTimes

	lines, err := ps.ReadLines("/proc/stat")
	if err != nil {
		fmt.Println("readlines", err.Error())
		return ret
	}

	if len(lines) > 0 && strings.HasPrefix(lines[0], "cpu") {
		fmt.Println("lines0", lines[0])
		ret, err = lineToCpuTimes(lines[0], true)
		if err != nil {
			fmt.Println("linetocpu", err.Error())
			return ret
		}
	}

	return ret
}

func (ps *PSUtils) PerCpuTimes() []CpuTimes {
	cpus := []CpuTimes{}
	lines, err := ps.ReadLines("/proc/stat")
	if err != nil {
		return cpus
	}

	for _, line := range lines {
		ct, err := lineToCpuTimes(line, false)
		if err != nil {
			continue
		}
		cpus = append(cpus, ct)
	}

	return cpus
}

func lineToCpuTimes(line string, total bool) (CpuTimes, error) {

	var ret CpuTimes
	sp := SplitString(line)

	if len(sp) < 8 {
		return ret, errors.New("less than 8 part, not related to cpu")
	} else if total == false && sp[0] != "cpu" && strings.HasPrefix(sp[0], "cpu") {
		ret = listSpToCpuTimes(sp)
		return ret, nil
	} else if total == true && sp[0] == "cpu" {
		ret = listSpToCpuTimes(sp)
		return ret, nil
	} else {
		return ret, errors.New(" not related to cpu")
	}
}

func listSpToCpuTimes(sp []string) CpuTimes {
	ret := CpuTimes{}
	if len(sp) >= 8 {
		ret.User = strToInt64(sp[1], 0)
		ret.Nice = strToInt64(sp[2], 0)
		ret.System = strToInt64(sp[3], 0)
		ret.Idle = strToInt64(sp[4], 0)
		ret.Iowait = strToInt64(sp[5], 0)
		ret.Irq = strToInt64(sp[6], 0)
		ret.Softirq = strToInt64(sp[7], 0)
	}

	if len(sp) >= 8 {
		ret.Steal = strToInt64(sp[8], 0)
	}

	if len(sp) >= 9 {
		ret.Guest = strToInt64(sp[9], 0)
	}

	if len(sp) >= 10 {
		ret.GuestNice = strToInt64(sp[10], 0)

	}
	return ret
}

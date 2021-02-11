package PSUtils

import (
	"fmt"
	"strings"
	"time"
)

var (
	GetBlocksCmd         = "ls /sys/block"
	CatProcDiskstatsCmd  = "cat /proc/diskstats"
	OmitDiskNamePrefixes = []string{"loop", "ram", "sr", "md", "dm-"}
)

type DiskOverallStats struct {
	ProcDiskStats
	ReadIOPS  int64
	WriteIOPS int64
	ReadSpeed int64
}

type DiskStats struct {
	Disks []struct {
		Path        string  `json:"path"`
		Fstype      string  `json:"fstype"`
		Total       int64   `json:"total"`
		Free        int64   `json:"free"`
		Used        int64   `json:"used"`
		UsedPercent float64 `json:"usedPercent"`
	} `json:"disks"`
}

/*
	1: https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
	 1  major number
	 2  minor mumber
	 3  device name
	 4  reads completed successfully
	 5  reads merged
	 6  sectors read
	 7  time spent reading (ms)
	 8  writes completed
	 9  writes merged
	10  sectors written
	11  time spent writing (ms)
	12  I/Os currently in progress
	13  time spent doing I/Os (ms)
	14  weighted time spent doing I/Os (ms)
*/
type ProcDiskStats struct {
	Major                       int
	Minor                       int
	DevName                     string
	ReadsCompletedSuccess       int64
	ReadsMerged                 int64
	SectorsRead                 int64
	TimeSpentReadingMS          int64
	WritesCompleted             int64
	WritesMerged                int64
	SectorsWritten              int64
	TimeSpentWritingMS          int64
	IOsCurrentlyInProgress      int64
	TimeSpentDoingIOsMS         int64
	WeightedTimeSpentDoingIOsMS int64

	ReadIOPS   int64
	WriteIOPS  int64
	ReadSpeed  int64 // bytes/s
	WriteSpeed int64 // bytes/s
}

func (ps *PSUtils) GetDiskOverallStats() *ProcDiskStats {
	cur := ps.GetSumProcDiskStats()
	curTM := time.Now().Unix()
	if ps.ProcDiskstatTmstamp != 0 {
		gap := curTM - ps.ProcDiskstatTmstamp
		if gap > 0 {
			if ps.LastDiskStat.ReadsCompletedSuccess > 0 {
				cur.ReadIOPS = (cur.ReadsCompletedSuccess - ps.LastDiskStat.ReadsCompletedSuccess) / gap
			}
			if ps.LastDiskStat.WritesCompleted > 0 {
				cur.WriteIOPS = (cur.WritesCompleted - ps.LastDiskStat.WritesCompleted) / gap
			}
			if ps.LastDiskStat.SectorsRead > 0 {
				cur.ReadSpeed = (cur.SectorsRead - ps.LastDiskStat.SectorsRead) * 512 / gap
			}
			if ps.LastDiskStat.SectorsWritten > 0 {
				cur.WriteSpeed = (cur.SectorsRead - ps.LastDiskStat.SectorsWritten) * 512 / gap
			}
		}
	}
	ps.LastDiskStat = cur
	ps.ProcDiskstatTmstamp = curTM

	return &cur
}

func (ps *PSUtils) GetSumProcDiskStats() ProcDiskStats {
	ret := ProcDiskStats{}
	devs, err := ps.GetNotVirtualBlockDeviceNames()
	if err != nil {
		return ret
	}

	stats := ps.parseProcDiskStats()
	for _, stat := range stats {
		if !ListContain(devs, stat.DevName) {
			continue
		}
		ret.ReadsCompletedSuccess += stat.ReadsCompletedSuccess
		ret.ReadsMerged += stat.ReadsMerged
		ret.SectorsRead += stat.SectorsRead
		ret.TimeSpentReadingMS += stat.TimeSpentReadingMS
		ret.WritesCompleted += stat.WritesCompleted
		ret.WritesMerged += stat.WritesMerged
		ret.SectorsWritten += stat.SectorsWritten
		ret.TimeSpentWritingMS += stat.TimeSpentWritingMS
		ret.IOsCurrentlyInProgress += stat.IOsCurrentlyInProgress
		ret.TimeSpentDoingIOsMS += stat.TimeSpentDoingIOsMS
		ret.WeightedTimeSpentDoingIOsMS += stat.WeightedTimeSpentDoingIOsMS
	}

	return ret
}

func (ps *PSUtils) parseProcDiskStats() []ProcDiskStats {
	ret := []ProcDiskStats{}
	s, err := ps.Exec(CatProcDiskstatsCmd)
	if err != nil {
		return nil
	}
	lines := SplitStringToLines(s)
	for _, line := range lines {
		fs := SplitString(line)
		if len(fs) < 14 {
			fmt.Println("fields in line of /proc/diskstats less than 14, not valid")
			continue
		}
		fmt.Printf("/proc/diskstats line: '%s'", line)
		ret = append(ret, ProcDiskStats{
			Major:                       strToInt(fs[0], 0),
			Minor:                       strToInt(fs[1], 0),
			DevName:                     fs[2],
			ReadsCompletedSuccess:       strToInt64(fs[3], 0),
			ReadsMerged:                 strToInt64(fs[4], 0),
			SectorsRead:                 strToInt64(fs[5], 0),
			TimeSpentReadingMS:          strToInt64(fs[6], 0),
			WritesCompleted:             strToInt64(fs[7], 0),
			WritesMerged:                strToInt64(fs[8], 0),
			SectorsWritten:              strToInt64(fs[9], 0),
			TimeSpentWritingMS:          strToInt64(fs[10], 0),
			IOsCurrentlyInProgress:      strToInt64(fs[11], 0),
			TimeSpentDoingIOsMS:         strToInt64(fs[12], 0),
			WeightedTimeSpentDoingIOsMS: strToInt64(fs[13], 0),
		})
	}
	return ret
}

func (ps *PSUtils) GetNotVirtualBlockDeviceNames() ([]string, error) {
	if ps.StorageDeviceNames != nil && len(ps.StorageDeviceNames) > 0 {
		return ps.StorageDeviceNames, nil
	}
	names, err := ps.getBlockDeviceNames()
	if err != nil {
		return nil, err
	}
	devs := FilterNonStorageDevice(names)
	ps.StorageDeviceNames = devs
	return devs, nil
}

func (ps *PSUtils) getBlockDeviceNames() ([]string, error) {
	s, err := ps.Exec(GetBlocksCmd)
	if err != nil {
		return nil, err
	}

	return SplitString(s), nil
}

func FilterNonStorageDevice(s []string) []string {
	// https://github.com/salesforce/LinuxTelemetry/blob/master/plugins/diskstats.py#L154
	var ret []string
	for _, i := range s {
		flag := true
		for _, name := range OmitDiskNamePrefixes {
			if strings.HasPrefix(i, name) {
				flag = false
			}
		}

		if flag {
			ret = append(ret, i)
		}
	}
	return ret
}

func (ps *PSUtils) TestDisk() {
	devs, err := ps.GetNotVirtualBlockDeviceNames()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, dev := range devs {
		fmt.Printf("line: '%s'\n", dev)
	}

	ps.parseProcDiskStats()
	// ps.GetDiskOverallStats()
}

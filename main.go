package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"madan.asia/ligulfzhou/gopsutil-mobile/PSUtils"
)

func main() {

	var err error
	fmt.Println("---------------Ubuntu--------------")
	// psutils := PSUtils.NewPSUtils("zhou", "", "192.168.8.194", "/Users/ligangzhou/.ssh/id_rsa", 22)
	psutils := PSUtils.NewPSUtils("root", "", "47.99.50.83", "/Users/ligangzhou/.ssh/id_rsa", "", 22)
	success, err := psutils.Connect()
	fmt.Printf("connection status: %t \n", success)
	if err != nil {
		log.Fatal(err)
	}

	s, err := psutils.Exec("dmesg | grep 'hostname'")
	if err == nil && PSUtils.StripString(s) != "" {
		r, _ := regexp.Compile(`Set hostname to <(?P<hostname>\S*)>\.`)
		m := r.FindStringSubmatch(s)
		if len(m) > 0 {
			fmt.Printf("'%s' \n", m[len(m)-1])
		}
	}

	dd, _ := psutils.ListDirectorys("/proc")
	fmt.Println("======list directory of /proc=====")
	for _, d := range dd {
		fmt.Printf("'%s'\n", d)
	}
	cnt := psutils.NumProcs()
	fmt.Printf("Number of Procs '%d'", cnt)
	fmt.Println("=============host=============")
	fmt.Println(psutils.PlatformInformation())
	fmt.Println("=============host=============")
	fmt.Println("=============mem=============")
	fmt.Println(psutils.VirtualMemory())
	fmt.Println("=============mem=============")
	fmt.Println("=============load=============")
	fmt.Println(psutils.ArgLoad())
	fmt.Println("=============load=============")
	fmt.Println("=============cpu count=============")
	fmt.Println(psutils.CPUCount(true))
	fmt.Println(psutils.CPUCount(false))
	fmt.Println(psutils.CPUInfo())
	err = psutils.GetMainInterface()
	fmt.Println(err)
	fmt.Println(psutils.NetworkInterface)
	fmt.Println(psutils.GetNetStats())
	fmt.Println("=============load=============")
	timer1 := time.NewTimer(1 * time.Second)
	<-timer1.C
	fmt.Println(psutils.GetNetStats())

	psutils.TestKernel()
	psutils.TestDisk()
	return

	fmt.Println("---------------Centos--------------")
	psutils = PSUtils.NewPSUtils("root", "", "192.168.8.193", "/Users/ligangzhou/.ssh/id_rsa", "", 22)
	success, err = psutils.Connect()
	fmt.Printf("connection status: %t \n", success)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=============host=============")
	fmt.Println(psutils.PlatformInformation())
	fmt.Println("=============host=============")
	fmt.Println("=============mem=============")
	fmt.Println(psutils.VirtualMemory())
	fmt.Println("=============mem=============")
	fmt.Println("=============load=============")
	fmt.Println(psutils.ArgLoad())
	fmt.Println("=============load=============")
	fmt.Println("=============cpu count=============")
	fmt.Println(psutils.CPUCount(true))
	// fmt.Println(psutils.CPUCount(false))
	fmt.Println(psutils.CPUInfo())
	fmt.Println("=============load=============")

	fmt.Println("---------------Debian--------------")
	psutils = PSUtils.NewPSUtils("root", "zhouligang", "192.168.8.135", "", "", 22)
	success, err = psutils.Connect()
	fmt.Printf("connection status: %t \n", success)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=============host=============")
	fmt.Println(psutils.PlatformInformation())
	fmt.Println("=============host=============")
	fmt.Println("=============mem=============")
	fmt.Println(psutils.VirtualMemory())
	fmt.Println("=============mem=============")
	fmt.Println("=============load=============")
	fmt.Println(psutils.ArgLoad())
	fmt.Println("=============load=============")
	fmt.Println("=============cpu count=============")
	fmt.Println(psutils.CPUCount(true))
	// fmt.Println(psutils.CPUCount(false))
	fmt.Println(psutils.CPUInfo())
	fmt.Println("=============load=============")

}

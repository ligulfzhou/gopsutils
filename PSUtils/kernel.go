package PSUtils

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	KERNEL_UNAME_CMD       = "uname -r"
	KERNEL_HOSTNAMECTL_CMD = "hostnamectl"
	KERNEL_PROC_CMD        = "cat /proc/version"
	KERNEL_DMSG_CMD        = "dmesg | grep 'Linux version'"

	KERNEL_ARCH_UNAME_CMD = "uname -m"
)

func (ps *PSUtils) GetKernalArch() string {
	s, err := ps.Exec(KERNEL_ARCH_UNAME_CMD)
	if err != nil {
		return ""
	}

	return StripString(s)
}

func (ps *PSUtils) GetKernelVersion() string {
	var kernel string
	var err error
	kernel, err = ps.getKernelVersionFromUnameCmd()
	if err == nil {
		return kernel
	}
	fmt.Printf("Uname failed: %s", err.Error())

	kernel, err = ps.getKernelVersionFromHostnamectlCmd()
	if err == nil {
		return kernel
	}
	fmt.Printf("Hostnamectl failed: %s", err.Error())

	kernel, err = ps.getKernelVersionFromProcCmd()
	if err == nil {
		return kernel
	}
	fmt.Printf("cat /proc/version failed: %s", err.Error())

	kernel, err = ps.getKernelVersionFromDmsgCmd()
	if err == nil {
		return kernel
	}
	fmt.Printf("cat dmsg failed: %s", err.Error())
	fmt.Println("get KernelVersion totally failed......though It may not happen, It still occurs")
	return ""
}

func (ps *PSUtils) getKernelVersionFromUnameCmd() (string, error) {
	/*
		4.18.0-193.19.1.el8_2.x86_64
	*/
	s, err := ps.Exec(KERNEL_UNAME_CMD)
	if err != nil {
		return "", err
	}
	return StripString(s), nil
}

func (ps *PSUtils) getKernelVersionFromHostnamectlCmd() (string, error) {
	/*
		   Static hostname: localhost.localdomain
		Transient hostname: EC-F4-BB-E3-41-F0
		         Icon name: computer-server
		           Chassis: server
		        Machine ID: 8ccb44cb4fa9440297cdd848099e75bf
		           Boot ID: 1ae8296a75a145599e9b38255e18f6a6
		  Operating System: CentOS Linux 8 (Core)
		       CPE OS Name: cpe:/o:centos:centos:8
		            Kernel: Linux 4.18.0-193.19.1.el8_2.x86_64
		      Architecture: x86-64
	*/
	s, err := ps.Exec(KERNEL_HOSTNAMECTL_CMD)
	if err != nil {
		return "", err
	}

	v := GetValueFromMapString(s, ":", "Kernel")
	return StripString(strings.TrimLeft(v, "Linux")), nil
	// return fmt.Sprintf("'%s'", StripString(strings.TrimLeft(v, "Linux"))), nil
}

func (ps *PSUtils) getKernelVersionFromProcCmd() (string, error) {
	/*
		Linux version 4.18.0-193.19.1.el8_2.x86_64 (mockbuild@kbuilder.bsys.centos.org) (gcc version 8.3.1 20191121 (Red Hat 8.3.1-5) (GCC)) #1 SMP Mon Sep 14 14:37:00 UTC 2020
	*/
	// s, err := ps.Exec(KERNEL_PROC_CMD)
	s, err := ps.Exec(KERNEL_PROC_CMD)
	if err != nil {
		return "", err
	}

	r, _ := regexp.Compile(`Linux version (?P<version>\S*) \(`)
	m := r.FindStringSubmatch(s)
	if len(m) > 0 {
		return m[len(m)-1], nil
	}
	return "", nil
}

func (ps *PSUtils) getKernelVersionFromDmsgCmd() (string, error) {
	/*
		[    0.000000] Linux version 4.4.0-166-generic (buildd@lcy01-amd64-020) (gcc version 5.4.0 20160609 (Ubuntu 5.4.0-6ubuntu1~16.04.10) ) #195-Ubuntu SMP Tue Oct 1 09:35:25 UTC 2019 (Ubuntu 4.4.0-166.195-generic 4.4.194)
	*/
	s, err := ps.Exec(KERNEL_DMSG_CMD)
	if err != nil {
		return "", err
	}

	r, _ := regexp.Compile(`Linux version (?P<version>\S*) \(`)
	m := r.FindStringSubmatch(s)
	if len(m) > 0 {
		return m[len(m)-1], nil
	}

	return "", nil
}

func (ps *PSUtils) TestKernel() {
	var kernel string
	var err error
	kernel, err = ps.getKernelVersionFromUnameCmd()
	fmt.Println(err)
	fmt.Printf("Uname: %s", kernel)

	kernel, err = ps.getKernelVersionFromHostnamectlCmd()
	fmt.Println(err)
	fmt.Printf("Hostname : %s", kernel)

	kernel, err = ps.getKernelVersionFromProcCmd()
	fmt.Println(err)
	fmt.Printf("cat /proc/version failed: %s", kernel)

	kernel, err = ps.getKernelVersionFromDmsgCmd()
	fmt.Println(err)
	fmt.Printf("cat dmsg failed: %s", kernel)
	fmt.Println("get KernelVersion totally failed......though It may not happen, It still occurs")
}

package main

import (
	"fmt"
	"log"
)

func main() {
	var err error
	fmt.Println("---------------Ubuntu--------------")
	psutils := NewPSUtils("zhou", "", "192.168.8.194", "/Users/ligangzhou/.ssh/id_rsa", 22, nil)
	success, err := psutils.Connect()
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

	fmt.Println("---------------Centos--------------")
	psutils = NewPSUtils("root", "", "192.168.8.193", "/Users/ligangzhou/.ssh/id_rsa", 22, nil)
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

	fmt.Println("---------------Debian--------------")
	psutils = NewPSUtils("root", "zhouligang", "192.168.8.135", "", 22, nil)
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
}

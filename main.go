package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	var err error
	psutils := NewPSUtils("zhou", "", "192.168.8.194", "/Users/ligangzhou/.ssh/id_rsa", 22, nil)
	success, err := psutils.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(success)
	exists := psutils.FileExists("/home/zhou/.ssh/known_hostss")
	fmt.Printf("/root/.ssh/known_hosts fileexists: %t\n", exists)

	exists = psutils.FileExists("/home/zhou/.ssh/known_hosts")
	fmt.Printf("/root/.ssh/known_hostss fileexists: %t\n", exists)


	psutils.PlatformInformationWithContext(context.Background())
}

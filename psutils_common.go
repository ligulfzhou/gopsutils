package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
)

// not implemented
func (ps *PSUtils) SudoExec(cmd string) (string, error){
	return "not implemented", nil
}

func (ps *PSUtils) Exec(cmd string) (string, error){

	session, err := ps.client.NewSession()
	if err != nil {
		return "", err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return "", err
	}

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		return "", err
	}

	fmt.Println("\n\n\n---------------exec cmd----------------")
	fmt.Printf("exec cmd: %s, res: %s \n", cmd, b.String())
	fmt.Println("---------------exec cmd---------------- \n\n\n ")
	return b.String(), nil
}

func (ps *PSUtils) FileExists(filename string) bool {
	_, err := ps.Exec(fmt.Sprintf("stat %s", filename))
	if err != nil {
		return false
	}

	return true
}

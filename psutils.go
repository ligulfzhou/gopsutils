package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"time"
)

type PSUtils struct {
	user, password,host, key string
	port int
	cipherList []string

	platform string

	client *ssh.Client
}

func NewPSUtils(user, password, host, key string, port int, cipherList []string) *PSUtils {

	return &PSUtils{
		user:       user,
		password:   password,
		host:       host,
		key:        key,
		port:       port,
		cipherList: cipherList,
		client:     nil,
	}
}

func (ps *PSUtils) Connect() (bool, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	if ps.key == "" {
		auth = append(auth, ssh.Password(ps.password))
	} else {
		pemBytes, err := ioutil.ReadFile(ps.key)
		if err != nil {
			return false, err
		}

		var signer ssh.Signer
		if ps.password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(ps.password))
		}
		if err != nil {
			return false, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	if len(ps.cipherList) == 0 {
		config = ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		}
	} else {
		config = ssh.Config{
			Ciphers: ps.cipherList,
		}
	}

	clientConfig = &ssh.ClientConfig{
		User:    ps.user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr = fmt.Sprintf("%s:%d", ps.host, ps.port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return false, err
	}
	ps.client = client
	return true, nil
}

func (ps *PSUtils) checkConn() {
	select {

	}
}
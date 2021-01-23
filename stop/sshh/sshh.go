package sshh

import (
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
)

// SSH Helper functions

func keyAuthMethod(keyPath string) ssh.AuthMethod {
	pemBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil
	}

	if !x509.IsEncryptedPEMBlock(block) {
		signer, _ := ssh.ParsePrivateKey(pemBytes)
		if signer == nil {
			return nil
		}
		return ssh.PublicKeys(signer)
	}
	return nil
}

func Connect(user string, addr string, keyPath string) *ssh.Client {
	keyAuth := keyAuthMethod(keyPath)
	if keyAuth == nil {
		return nil
	}

	client, _ := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{keyAuth},
		HostKeyCallback: func(_ string, _ net.Addr, _ ssh.PublicKey) error {
			return nil
		},
	})

	return client
}

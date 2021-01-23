package remote

import (
	"bytes"
	"fmt"
	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"os"
	"os/user"
	"path/filepath"
	"stop/sshh"
	"strconv"
)

type Client struct {
	*ssh.Client
	Alias string
	Addr  string
}

func NewClient(alias string) (client *Client) {
	currentUser, err := user.Current()
	if err != nil {
		return
	}
	userName := ssh_config.Get(alias, "User")
	if userName == "" {
		userName = currentUser.Username
	}

	identityFilePath := ssh_config.Get(alias, "IdentityFilePath")
	if identityFilePath == "" {
		identityFilePath = filepath.Join(currentUser.HomeDir, ".ssh", "id_rsa")
	}
	if _, err = os.Stat(identityFilePath); err != nil {
		return
	}

	host := ssh_config.Get(alias, "HostName")
	port := ssh_config.Get(alias, "Port")
	if host == "" || port == "" {
		return
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	c := sshh.Connect(userName, addr, identityFilePath)
	if c != nil {
		return &Client{Client: c, Alias: alias, Addr: addr}
	}

	return
}

func (c Client) Close() {
	_ = c.Client.Close()
}

func (c Client) RunCommand(cmd string) (stdout string, err error) {
	session, err := c.Client.NewSession()
	if err != nil {
		return
	}
	defer func() { _ = session.Close() }()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(cmd)
	if err != nil {
		return
	}
	stdout = string(buf.Bytes())
	return
}

func (c Client) FloatValue(cmd string) (float64, error) {
	output, err := c.RunCommand(prepCmd(cmd))
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(output, 64)
}

func (c Client) UintValue(cmd string) (uint64, error) {
	output, err := c.RunCommand(prepCmd(cmd))
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(output, 10, 64)
}

func (c Client) StringValue(cmd string) (string, error) {
	return c.RunCommand(prepCmd(cmd))
}

// prepCmd trim spaces and return the first line without trailing newline.
func prepCmd(cmd string) string {
	return fmt.Sprintf("%s | sed -e 's/^ *//g' -e 's/ *$//g' | sed -n '1 p' | tr -d '\n'", cmd)
}

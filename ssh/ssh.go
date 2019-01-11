package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	Host       string
	User       string
	Password   string
	PrivateKey string
	*ssh.Client
}

func (this *Client) Connect() error {
	conf := ssh.ClientConfig{
		User:            this.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if this.Password != "" {
		conf.Auth = append(conf.Auth, ssh.Password(this.Password))
	}

	if privateKey, err := getPrivateKey(this.PrivateKey); err == nil {
		conf.Auth = append(conf.Auth, privateKey)
	}

	client, err := ssh.Dial("tcp", this.Host, &conf)
	if err != nil {
		return fmt.Errorf("unable to connect: %v", err)
	}

	this.Client = client

	return nil
}

// Close the connection
func (this *Client) Close() {
	this.Client.Close()
}

// Get the private key for current user
func getPrivateKey(privateKey string) (ssh.AuthMethod, error) {

	key := []byte(privateKey)

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %v", err)
	}

	return ssh.PublicKeys(signer), nil
}

func CreateTerminalModes() *ssh.TerminalModes {
	return &ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
}

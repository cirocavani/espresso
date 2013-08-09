package main

import (
	"bytes"
	"flag"
	"fmt"
	"runtime"
)

import "code.google.com/p/go.crypto/ssh"

var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")
var optUsername = flag.String("username", "", "Username")
var optPassword = flag.String("password", "", "Password")
var optServer = flag.String("server", "", "SSH server")

type password string

func (this password) Password(user string) (password string, err error) {
	password = string(this)
	err = nil
	return
}

func main() {
	fmt.Println("SSH Client")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	config := &ssh.ClientConfig{
		User: *optUsername,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthPassword(password(*optPassword)),
		},
	}
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", *optServer, config)
	if err != nil {
		panic("Unable to connect: " + err.Error())
	}
	defer conn.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := conn.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println("Output:")
	fmt.Print(b.String())
	fmt.Println("...")

	fmt.Println("bye!")
}

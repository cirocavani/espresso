SSH Package:

http://godoc.org/code.google.com/p/go.crypto/ssh

	(Mercurial is required)
	export GOPATH=<SOURCE FOLDER>/espresso/lib
	go get code.google.com/p/go.crypto/ssh

Output:

	go run sample/ssh.go -username cavani -password cavani -server 127.0.0.1:22
	SSH Client
	Threads: 8
	Output:
	cavani
	...
	bye!

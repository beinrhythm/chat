package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	_connsMu sync.RWMutex

	_conns = make(map[net.Conn]struct{})
)

// Call net.Listen from the main() func.
func main() {
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		// handle the error, e.g. `log.Fatal(err)`
		log.Fatal(err)
	}
	fmt.Println("Listening on ", ln.Addr())
	for {
		c, err := ln.Accept()
		if err == nil {
			// do something with `c`
			fmt.Println("Connection: ", c)
			// Start goroutines by prepending the `go` keyword to call the serve function
			go serve(c)
		}
	}
}

func serve(c net.Conn) {
	register(c, true)
	defer register(c, false)
	// handle the connection
	fmt.Fprintf(c, "Hello, there %v\n", c.RemoteAddr())

	// After registration
	bs := bufio.NewScanner(c)
	for bs.Scan() {
		// bs.Text() will contain the most recent line
		msg := bs.Text()
		_connsMu.RLock()
		for peer := range _conns {
			if peer != c {
				fmt.Fprintf(peer, "%v: %s\n", c.RemoteAddr(), msg)
			}
		}
		_connsMu.RUnlock()
	}
	if err := bs.Err(); err != nil {
		// bs.Err() returns any error encountered during the scan.
		log.Println("error reading from the connection: ", err)
	}
}

func register(c net.Conn, ok bool) {
	_connsMu.Lock()
	defer _connsMu.Unlock()

	if ok {
		fmt.Printf("Client connected: %v\n", c.RemoteAddr())
		_conns[c] = struct{}{}
	} else {
		fmt.Printf("Client disconnected: %v\n", c.RemoteAddr())
		delete(_conns, c)
	}
}

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"

	"github.com/humbertovnavarro/tcpesame/pkg/facade"
)

var (
	maxConns    int64 = 1000
	activeConns int64
)

func main() {
	for port := 1000; port <= 30000; port++ {
		go func(p int) {
			addr := ":" + strconv.Itoa(p)
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return
			}
			fmt.Printf("[+] Listening on %s\n", addr)
			for {
				conn, err := ln.Accept()
				if err != nil {
					return
				}
				current := atomic.AddInt64(&activeConns, 1)
				if current > maxConns {
					atomic.AddInt64(&activeConns, -1)
					fmt.Printf("[-] Connection dropped (limit reached: %d)\n", maxConns)
					conn.Close()
					continue
				}
				go func() {
					handleConn(conn, port)
					conn.Close()
				}()
			}
		}(port)
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	fmt.Println("\n[!] Caught signal — shutting down gracefully...")
	fmt.Printf("[*] Active connections at shutdown: %d\n", atomic.LoadInt64(&activeConns))
}

func handleConn(conn net.Conn, port int) {
	switch port {
	case 22:
		facade.SSHFacade(conn)
	case 2222:
		facade.SSHFacade(conn)
	default:
		facade.HTTPFacade(conn)
	}
}

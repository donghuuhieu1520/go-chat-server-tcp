package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var (
	conns   []net.Conn
	connCh  = make(chan net.Conn)
	closeCh = make(chan net.Conn)
	msgCh   = make(chan string)
	errCh   = make(chan error)
)

func main() {
	server, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	fmt.Println("Server is running on port 3000")
	go handleRequest(server)

	for {
		select {
		case conn := <-connCh:
			go onMessage(conn)
		case conn := <-closeCh:
			fmt.Println("A client has exited")
			removeConn(conn)
		case msg := <-msgCh:
			fmt.Println(msg)
		case err := <-errCh:
			log.Fatal(err)
		}
	}
}

func handleRequest(server net.Listener) {
	for {
		conn, err := server.Accept()
		if err != nil {
			errCh <- err
			break
		}

		conns = append(conns, conn)
		connCh <- conn
	}
}

func onMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msgCh <- msg
		publishMessage(conn, msg)
	}
	closeCh <- conn
}

func removeConn(conn net.Conn) {
	var i int
	for i = range conns {
		if conn == conns[i] {
			break
		}
	}
	conns = append(conns[:i], conns[i+1:]...)
}

func publishMessage(conn net.Conn, msg string) {
	for _, c := range conns {
		if c != conn {
			c.Write([]byte(msg))
		}
	}
}

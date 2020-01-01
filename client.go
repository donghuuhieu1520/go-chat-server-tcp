package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Print("What is your name: ")
	nameReader := bufio.NewReader(os.Stdin)
	name, _ := nameReader.ReadString('\n')
	name = name[:len(name)-1]

	for {
		msgReader := bufio.NewReader(os.Stdin)
		msg, _ := msgReader.ReadString('\n')
		msg = fmt.Sprintf("%s: %s\n", name, msg[:len(msg)-1])
		conn.Write([]byte(msg))
	}
}

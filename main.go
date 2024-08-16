package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	mode := flag.String("mode", "server", "Mode is set to server or client")

	flag.Parse()

	if mode == nil || *mode == "server" {
		server()
	} else {
		client()
	}

	fmt.Println("")
}

func client() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Text to send: ")

		input, _ := reader.ReadString('\n')
		conn.Write([]byte(input))
		message, _ := bufio.NewReader(conn).ReadString('\n')

		log.Print("Server relay:", message)
	}
}

func server() {
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		return
	}

	log.Println("Client message:", string(buffer[:len(buffer)-1]))

	conn.Write(buffer)

	handleConnection(conn)
}

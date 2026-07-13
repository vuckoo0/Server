package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	buffer = make([]byte, 1024)
)

func handleConn(conn net.Conn) {

	fmt.Printf("[+] Accepted a connection form %v\n", conn.LocalAddr())

	defer conn.Close()

	for {

		n, err := conn.Read(buffer)
		if err != nil {

			if errors.Is(err, io.EOF) {

				fmt.Printf("[-] Client %v disconected\n", conn.LocalAddr())

			} else {

				fmt.Printf("[-] Error in reading from the client %v: %e", conn.LocalAddr(), err)
			}
		}

		message := string(buffer[:n])
		log.Printf("[%v] %v\n", conn.LocalAddr(), message)

		if message == "exit()" {
			log.Printf("[.] Client is disconecting...")
		}

		_, err = conn.Write([]byte("ok"))
		if err != nil {
			log.Fatal("Error in sending message to the server: ", err)
		}
	}
}

func main() {

	server, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal("[-] Server didnt start properly")
	}

	log.Default()
	fmt.Println("[+] Oppend server on port 8080...")

	for {

		conn, err := server.Accept()

		if err != nil {
			continue
		}

		go handleConn(conn)
	}
}

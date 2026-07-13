package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	buffer      = make([]byte, 1024)
	activeUsers = map[net.Addr]string{}
)

func handleConn(conn net.Conn) {

	defer conn.Close()

	for {

		n, err := conn.Read(buffer)
		if err != nil {

			if errors.Is(err, io.EOF) {

				fmt.Printf("[-][%s] Client unexpectedly disconected\n", activeUsers[conn.RemoteAddr()])

			} else {

				fmt.Printf("[-][%s] Error in reading from the client %v: %e", activeUsers[conn.RemoteAddr()], conn.RemoteAddr(), err)
			}
		}

		message := string(buffer[:n])
		log.Printf("[%s] %v\n", activeUsers[conn.RemoteAddr()], message)

		if message == "exit()" {
			log.Printf("[.][%s] Client is disconecting...", activeUsers[conn.RemoteAddr()])
			break
		}

		_, err = conn.Write([]byte("ok"))
		if err != nil {
			log.Fatal("[-] Error in sending message to the server: ", err)
		}
	}
}

func main() {

	server, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal("[-] Server didnt start properly")
	}

	fmt.Println("[+] Oppend server on port 8080...")

	for {

		conn, err := server.Accept()

		if err != nil {
			continue
		}

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("[-] Error in username obtaning")
			continue
		}

		username := string(buffer[:n])
		activeUsers[conn.RemoteAddr()] = username

		fmt.Printf("[+][Addr: %v] Accepted a connection form %s\n", conn.RemoteAddr(), username)

		go handleConn(conn)
	}
}

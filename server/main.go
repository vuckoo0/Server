package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"main/server/recorder"
	"net"
	"time"
)

var (
	buffer        = make([]byte, 1024)
	activeUsers   = map[net.Addr]string{}
	messageChanel = make(chan recorder.Row)
)

func handleConn(conn net.Conn) {

	defer conn.Close()

	currentRow := new(recorder.Row)

	for {

		n, err := conn.Read(buffer)
		if err != nil {

			if errors.Is(err, io.EOF) {

				fmt.Printf("[-][%s] Client unexpectedly disconected\n", activeUsers[conn.RemoteAddr()])

			} else {

				fmt.Printf("[-][%s] Error in reading from the client %v: %e", activeUsers[conn.RemoteAddr()], conn.RemoteAddr(), err)
			}
		}

		currentRow.Message = string(buffer[:n])
		currentRow.Ip = conn.RemoteAddr().String()
		currentRow.User = activeUsers[conn.RemoteAddr()]
		currentRow.Time = time.Now().Format("2006-01-02 15:04:05")

		messageChanel <- *currentRow

		if currentRow.Message == "exit()" {
			log.Printf("[.][%s] Client is disconecting...", activeUsers[conn.RemoteAddr()])
			break
		}

		log.Printf("[%s] %v\n", activeUsers[conn.RemoteAddr()], currentRow.Message)

		_, err = conn.Write([]byte("ok"))
		if err != nil {
			log.Fatal("[-] Error in sending message to the server: ", err)
		}
	}
}

func main() {

	db, err := sql.Open("mysql", "root:Vucko1602!@tcp(localhost:3306)/first_db")
	if err != nil {
		log.Fatal("[-] Error in opening database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("[-] Error in pinging database: ", err)
	}

	go recorder.Recorder(messageChanel, db)

	server, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal("[-] Server didnt start properly")
	}

	fmt.Println("[+] Opened server on port 8080...")

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

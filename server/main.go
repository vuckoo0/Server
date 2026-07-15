package main

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"main/server/recorder"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
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

				log.Printf("[-][%s] Client unexpectedly disconected\n", activeUsers[conn.RemoteAddr()])

			} else {

				log.Printf("[-][%s] Error in reading from the client %v: %e", activeUsers[conn.RemoteAddr()], conn.RemoteAddr(), err)
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	dsn := os.Getenv("DSN")
	db, err := sql.Open("mysql", dsn)
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

	log.Println("[+] Opened server on port 8080...")

	for {

		conn, err := server.Accept()

		if err != nil {
			continue
		}

		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("[-] Error in username obtaning")
			continue
		}

		username := string(buffer[:n])
		activeUsers[conn.RemoteAddr()] = username

		log.Printf("[+][Addr: %v] Accepted a connection form %s\n", conn.RemoteAddr(), username)

		go handleConn(conn)
	}
}

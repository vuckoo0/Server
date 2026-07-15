package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strings"
)

var (
	reader = bufio.NewReader(os.Stdin)
	buffer = make([]byte, 1024)
)

func readLine(promt string) string {

	fmt.Print(promt)

	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func main() {

	conn, err := net.Dial("tcp", "192.168.1.237:8080")
	if err != nil {
		log.Fatal("[-] Error in connecting to the server")
	}

	defer conn.Close()

	username := readLine("[+] Enter the username that will be used to identify you on the server: ")

	_, err = conn.Write([]byte(username))
	if err != nil {
		fmt.Println("[-] Error in sending the username to the server")
	}

	for {

		message := readLine("[+] Enter a message for the server: ")

		if message == "exit()" {

			buffer = []byte("exit()")
			_, err := conn.Write(buffer)

			if err != nil {
				log.Fatal("[-] Error in sending message to the server")
			}

			log.Println("[.] Disconnecting from the seerver...")
			break
		}

		buffer = []byte(message)

		_, err := conn.Write(buffer)

		if err != nil {
			log.Fatal("[-] Error in sending message to the server")
		}

		n, err := conn.Read(buffer)

		if err != nil {
			log.Fatal("[-] Error in reading the message from the server")
		}

		if slices.Equal(buffer[:n], []byte("exit()")) {
			fmt.Println("[.] Exiting program...")
			break
		}

		fmt.Println(string(buffer[:n]))
	}
}

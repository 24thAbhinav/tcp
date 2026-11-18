package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var terminalMu sync.Mutex

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}
	username = strings.TrimSpace(username)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	serverReader := bufio.NewReader(conn)
	//sending the username first
	conn.Write([]byte(username + "\n"))

	// for reading the text sent from the server
	go func() {
		for {
			message, err := serverReader.ReadString('\n')
			if err != nil {
				return
			}

			terminalMu.Lock()
			fmt.Print("\r\033[K" + message)
			fmt.Printf("%s: ", username)
			terminalMu.Unlock()
		}
	}()

	for {

		terminalMu.Lock()
		fmt.Printf("%s: ", username)
		terminalMu.Unlock()

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input error:", err)
			return
		}

		_, err = conn.Write([]byte(username + ": " + message))
		if err != nil {
			fmt.Println("Send error:", err)
			return
		}
	}
}

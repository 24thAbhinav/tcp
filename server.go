package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var connectionMap = make(map[string]net.Conn)
var connectionMapMut sync.Mutex

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)

	username, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Issues receiving username")
		return
	}

	//removing the \n from the end of the username
	username = strings.TrimRight(username, "\r\n")

	fmt.Printf("%s joined the server\n", username)

	connectionMapMut.Lock()
	connectionMap[username] = c
	connectionMapMut.Unlock()

	broadcast(username, username+" joined the server")

	for {

		//receiving msgs from the clients in format username: text
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				broadcast(username, username+" left the server")
				connectionMapMut.Lock()
				delete(connectionMap, username)
				connectionMapMut.Unlock()
			} else {
				log.Println("Read error:", err)
			}
			return
		}

		// drop the framing newline; we add exactly one back when sending
		message = strings.TrimRight(message, "\r\n")

		// ignoring the username the client sent; use the one from the connection
		parts := strings.SplitN(message, ": ", 2)
		if len(parts) != 2 {
			continue
		}
		txt := parts[1]

		broadcast(username, username+": "+txt)

		fmt.Printf("%s: %s\n", username, txt)
	}
}

// sends text to every connected client except sender
func broadcast(sender, text string) {
	//take Screenshot of the map and then send
	connectionMapMut.Lock()

	clients := make([]net.Conn, 0, len(connectionMap))

	for usernm, conn := range connectionMap {
		if usernm != sender {
			clients = append(clients, conn)
		}
	}

	connectionMapMut.Unlock()

	for _, conn := range clients {
		if _, err := conn.Write([]byte(text + "\n")); err != nil {
			log.Println("Send error:", err)
		}
	}
}

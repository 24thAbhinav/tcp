package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("You: ")

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input error:", err)
			return
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Send error:", err)
			return
		}
	}
}

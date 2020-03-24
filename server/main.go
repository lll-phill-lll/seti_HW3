package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
		}

		go handleConnection(conn, false)
	}
}


func checkAuth(authStr string) (bool, error) {
	strings.Split(authStr)
}

func bytesToMap(bytes []byte) (map[string]string, error) {
	splitted := strings.Split(string(bytes), " ")
	if len(splitted) % 2 == 1 {
		return nil, errors.New("not all keys has values")

	}

}

func handleConnection(conn net.Conn, isAuth bool) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		log.Println("client left..", err)
		conn.Close()
		return
	}

	message := string(bufferBytes)

	if !isAuth {
		allow, err := checkAuth()
	}

	clientAddr := conn.RemoteAddr().String()
	response := fmt.Sprintf(message + " from " + clientAddr + "\n")

	log.Println(response)

	conn.Write([]byte("you sent: " + response))
	if message == "STOP\n" {
		log.Println("good bye")
		os.Exit(0)
	}

	handleConnection(conn, true)
}

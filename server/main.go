package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
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

		go handleAuth(conn)
	}
}

const (
	loginKey = "LOGIN"
	passwordKey = "PASSW"

	LOGIN_ERROR = "ERROR NO_LOGIN"
	PASSWORD_ERROR = "ERROR NO_PASSWORD"
	ARGUMENTS_NUMBER_ERROR = "ERROR ARGUMENTS_NUMBER"
	WRONG_LOGIN_OR_PASSWORD_ERROR = "ERROR LOGIN_OR_PASSWORD"

	SUCCESS_AUTH = "SUCC AUTH"

	ROOM = "ROOM"
	NEW = "NEW"
	GAME_ID_KEY = "GAME_ID"
)

func checkAuth(commands map[string]string) error {
	login, ok := commands[loginKey]
	if !ok {
		return errors.New(LOGIN_ERROR)
	}
	passw, ok := commands[passwordKey]
	if !ok {
		return errors.New(PASSWORD_ERROR)
	}

	if strings.Trim(login, "\n") == "yourmom" && strings.Trim(passw, "\n") == "gay" {
		return nil
	}

	return errors.New(WRONG_LOGIN_OR_PASSWORD_ERROR)
}

func bytesToMap(bytes []byte) (map[string]string, error) {
	splitted := strings.Split(string(bytes), " ")
	if len(splitted) % 2 == 1 {
		return nil, errors.New("not all keys has values")
	}

	// TODO add validation
	commands := make(map[string]string)
	key := ""
	for i := 0; i < len(splitted); i++ {
		if i % 2 == 0 {
			key = splitted[i]
		} else {
			commands[key] = splitted[i]
		}
	}

	return commands, nil
}

func handleAuth(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println("client left..", err)
		conn.Close()
		return
	}

	defer handleAuth(conn)

	clientAddr := conn.RemoteAddr().String()
	input := fmt.Sprintf(string(bufferBytes)+ " from " + clientAddr + "\n")
	log.Println(input)

	commands, err := bytesToMap(bufferBytes)
	if err != nil {
		conn.Write([]byte(err.Error() + "\n"))
		log.Println("error with command", err, string(bufferBytes))
		return
	}

	err = checkAuth(commands)
	if err != nil {
		conn.Write([]byte(err.Error() + "\n"))
		log.Println("auth fail", err, string(bufferBytes))
		return
	}

	log.Println("auth succ", err, string(bufferBytes))
	conn.Write([]byte(SUCCESS_AUTH + "\n"))
	handleCommands(conn)
}

func handleCommands(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		log.Println("client left..", err)
		conn.Close()
		return
	}

	defer handleCommands(conn)


	splitted := strings.Split(string(bufferBytes), " ")
	if len(splitted) != 2 {
		log.Println("Not 2 arguments")
		conn.Write([]byte(ARGUMENTS_NUMBER_ERROR + "\n"))
		return
	}
	splitted[1] = strings.Trim(splitted[1], "\n")

	if splitted[0] == ROOM && splitted[1] == NEW {
		conn.Write([]byte(GAME_ID_KEY + " 5\n"))
		// TODO: add generator
		// TODO: redirect to game handler
		return
	}
	if splitted[0] == ROOM {
		conn.Write([]byte(GAME_ID_KEY + " " + splitted[1] + "\n"))
		return
	}
}

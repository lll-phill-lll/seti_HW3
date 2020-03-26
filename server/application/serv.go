package application

import (
	"bufio"
	"chess/server/constants"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type Serv struct {
	Rand *Rand
}

// 127.0.0.1:8080
func (s *Serv) StartServe(address string) {
	listener, err := net.Listen("tcp", address)
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

func checkAuth(commands map[string]string) error {
	login, ok := commands[constants.LOGIN_KEY]
	if !ok {
		return errors.New(constants.LOGIN_ERROR)
	}
	passw, ok := commands[constants.PASSWORD_KEY]
	if !ok {
		return errors.New(constants.PASSWORD_ERROR)
	}

	if strings.Trim(login, "\n") == "yourmom" && strings.Trim(passw, "\n") == "gay" {
		return nil
	}

	return errors.New(constants.WRONG_LOGIN_OR_PASSWORD_ERROR)
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
	conn.Write([]byte(constants.SUCCESS_AUTH + "\n"))
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
		conn.Write([]byte(constants.ARGUMENTS_NUMBER_ERROR + "\n"))
		return
	}
	splitted[1] = strings.Trim(splitted[1], "\n")

	if splitted[0] == constants.ROOM && splitted[1] == constants.NEW {
		conn.Write([]byte(constants.GAME_ID_KEY + " 5\n"))
		// TODO: add generator
		// TODO: redirect to game handler
		return
	}
	if splitted[0] == constants.ROOM {
		conn.Write([]byte(constants.GAME_ID_KEY + " " + splitted[1] + "\n"))
		return
	}
	conn.Write([]byte(constants.UNKNOWN_COMMAND_ERROR + "\n"))
}

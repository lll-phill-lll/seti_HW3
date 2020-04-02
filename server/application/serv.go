package application

import (
	"bufio"
	"chess/server/constants"
	"chess/server/game"
	"chess/server/room"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Serv struct {
	Rand  *Rand
	Users []room.Player
	Rooms []room.Room
	mu    sync.Mutex
}

func (s *Serv) GetRoom(id int, player room.Player) (room.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for num, r := range s.Rooms {
		if r.Game.GetID() == id {
			if r.Has2Players {
				if r.WhitePlayer.Login == player.Login || r.BlackPlayer.Login == player.Login {
					return r, nil
				} else {
					return room.Room{}, errors.New(constants.INCORRECT_ROOM_ID_ERROR)
				}
			} else {
				s.Rooms[num].BlackPlayer = player
				s.Rooms[num].Has2Players = true
				return s.Rooms[num], nil
			}
		}
	}
	return room.Room{}, errors.New(constants.INCORRECT_ROOM_ID_ERROR)
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

		go s.handleAuth(conn)
	}

}

func (s *Serv) addNewRoom(room room.Room) {
	s.mu.Lock()
	defer s.mu.Unlock()
	MyLog.Println("New room", room)
	s.Rooms = append(s.Rooms, room)
}

func (s *Serv) checkAdmin(commands map[string]string) bool {
	login, ok := commands[constants.LOGIN_KEY]
	if !ok {
		return false
	}
	password, ok := commands[constants.PASSWORD_KEY]
	if !ok {
		return false
	}

	if login == constants.ADMIN_LOGIN && strings.Trim(password, "\n") == constants.ADMIN_PASSWORD {
		MyLog.Println("admin logged")
		return true
	}
	return false
}

func (s *Serv) checkAuth(commands map[string]string) (room.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	login, ok := commands[constants.LOGIN_KEY]
	if !ok {
		return room.Player{}, errors.New(constants.LOGIN_ERROR)
	}
	isLogin := false
	curUser := room.Player{}
	for _, user := range s.Users {
		if user.Login == strings.Trim(login, "\n") {
			isLogin = true
			curUser = user
			break
		}
	}
	if !isLogin {
		return room.Player{}, errors.New(constants.WRONG_LOGIN_OR_PASSWORD_ERROR)
	}

	passw, ok := commands[constants.PASSWORD_KEY]
	if !ok {
		return room.Player{}, errors.New(constants.PASSWORD_ERROR)
	}

	if strings.Trim(passw, "\n") != curUser.Password {
		return room.Player{}, errors.New(constants.WRONG_LOGIN_OR_PASSWORD_ERROR)
	}

	return curUser, nil
}

func (s *Serv) checkRegistration(commands map[string]string) (room.Player, bool, error) {
	_, ok := commands[constants.REG_KEY]
	if !ok {
		return room.Player{}, false, nil
	}

	login, ok := commands[constants.LOGIN_KEY]
	if !ok {
		return room.Player{}, true, errors.New(constants.LOGIN_ERROR)
	}

	password, ok := commands[constants.PASSWORD_KEY]
	if !ok {
		return room.Player{}, true, errors.New(constants.PASSWORD_ERROR)
	}

	curUser := room.Player{
		Games:    []game.Game{},
		Login:    login,
		Password: strings.Trim(password, "\n"),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range s.Users {
		if user.Login == login {
			return room.Player{}, true, errors.New(constants.NON_UNIQUE_LOGIN_ERROR)
		}
	}
	MyLog.Println("registred user", curUser)
	s.Users = append(s.Users, curUser)

	return curUser, true, nil
}

func (s *Serv) bytesToMap(bytes []byte) (map[string]string, error) {
	splitted := strings.Split(string(bytes), " ")
	if len(splitted)%2 == 1 {
		return nil, errors.New("not all keys has values")
	}

	// TODO add validation
	commands := make(map[string]string)
	key := ""
	for i := 0; i < len(splitted); i++ {
		if i%2 == 0 {
			key = splitted[i]
		} else {
			commands[key] = splitted[i]
		}
	}

	return commands, nil
}

func (s *Serv) handleAuth(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println("client left..", err)
		conn.Close()
		return
	}

	defer s.handleAuth(conn)

	clientAddr := conn.RemoteAddr().String()
	input := fmt.Sprintf(string(bufferBytes) + " from " + clientAddr + "\n")
	log.Println(input)

	commands, err := s.bytesToMap(bufferBytes)
	if err != nil {
		conn.Write([]byte(err.Error() + "\n"))
		log.Println("error with command", err, string(bufferBytes))
		return
	}

	player, isReg, err := s.checkRegistration(commands)
	if err != nil {
		conn.Write([]byte(err.Error() + "\n"))
		log.Println("auth fail", err, string(bufferBytes))
		return
	}

	if isReg {
		log.Println("reg succ", err, string(bufferBytes))
		conn.Write([]byte(constants.SUCCESS_REG + "\n"))
		s.handleCommands(conn, player)
		return
	}

	isAdmin := s.checkAdmin(commands)
	if isAdmin {
		conn.Write([]byte(constants.SUCCESS_AUTH + "\n"))
		s.handleAdmin(conn)
		return
	}

	player, err = s.checkAuth(commands)
	if err != nil {
		conn.Write([]byte(err.Error() + "\n"))
		log.Println("auth fail", err, string(bufferBytes))
		return
	}

	log.Println("auth succ", err, string(bufferBytes))
	MyLog.Println("auth success", player)
	conn.Write([]byte(constants.SUCCESS_AUTH + "\n"))
	s.handleCommands(conn, player)
}

func (s *Serv) handleAdmin(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		log.Println("client left..", err)
		conn.Close()
		return
	}

	defer s.handleAdmin(conn)

	splitted := strings.Split(string(bufferBytes), " ")
	if len(splitted) != 2 {
		log.Println("Not 2 arguments")
		conn.Write([]byte(constants.ARGUMENTS_NUMBER_ERROR + "\n"))
		return
	}
	splitted[1] = strings.Trim(splitted[1], "\n")

	if splitted[0] == constants.ADMIN_GET_KEY && splitted[1] == constants.ADMIN_PLAYERS {
		allUsers := ""
		for _, user := range s.Users {
			allUsers += user.Login + "\n"
		}
		conn.Write([]byte(allUsers + "\n"))
		return
	}

	conn.Write([]byte(constants.UNKNOWN_COMMAND_ERROR + "\n"))
}

func (s *Serv) handleCommands(conn net.Conn, player room.Player) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		log.Println("client left..", err)
		conn.Close()
		return
	}

	defer s.handleCommands(conn, player)
	MyLog.Println("New command from", player, "is:", string(bufferBytes))

	splitted := strings.Split(string(bufferBytes), " ")
	if len(splitted) != 2 {
		log.Println("Not 2 arguments")
		conn.Write([]byte(constants.ARGUMENTS_NUMBER_ERROR + "\n"))
		return
	}
	splitted[1] = strings.Trim(splitted[1], "\n")

	if splitted[0] == constants.ROOM && splitted[1] == constants.NEW {
		newRoomID := s.Rand.getNextID()
		conn.Write([]byte(fmt.Sprintf("%s %d\n", constants.GAME_ID_KEY, newRoomID)))
		newGame := game.NewGame(newRoomID)
		currRoom := room.Room{
			WhitePlayer: player,
			BlackPlayer: room.Player{},
			Game:        newGame,
			Has2Players: false,
		}
		s.addNewRoom(currRoom)
		s.handleRoom(conn, currRoom, player)
		return
	}
	if splitted[0] == constants.ROOM {
		id, err := strconv.Atoi(splitted[1])
		if err != nil {
			conn.Write([]byte(constants.INCORRECT_ROOM_ID_ERROR + "\n"))
			return
		}
		currRoom, err := s.GetRoom(id, player)
		if err != nil {
			conn.Write([]byte(constants.INCORRECT_ROOM_ID_ERROR + "\n"))
			return
		}
		conn.Write([]byte(constants.SUCCESS_ROOM + "\n"))
		s.handleRoom(conn, currRoom, player)
		return

	}
	conn.Write([]byte(constants.UNKNOWN_COMMAND_ERROR + "\n"))
}

func (s *Serv) handleRoom(conn net.Conn, currRoom room.Room, player room.Player) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		log.Println("client left..", err)
		conn.Close()
		return
	}

	defer s.handleRoom(conn, currRoom, player)
	MyLog.Println("COMMAND", string(bufferBytes), "from", player)

	splitted := strings.Split(string(bufferBytes), " ")
	if len(splitted) != 2 {
		log.Println("Not 2 arguments")
		conn.Write([]byte(constants.ARGUMENTS_NUMBER_ERROR + "\n"))
		return
	}

	command := splitted[0]
	if command != constants.MOVE_KEY {
		log.Println("incorrect move" + string(bufferBytes))
		conn.Write([]byte(constants.UNKNOWN_COMMAND_ERROR + "\n"))
		return
	}

	order := currRoom.Game.GetOrder()
	if currRoom.BlackPlayer.Login == player.Login {
		if order%2 != 1 {
			conn.Write([]byte(constants.NOT_YOUR_TURN_ERROR + "\n"))
			return
		}
	} else {
		if order%2 != 0 {
			conn.Write([]byte(constants.NOT_YOUR_TURN_ERROR + "\n"))
			return
		}
	}

	move := strings.Split(splitted[1], "-")

	if len(move) != 2 {
		log.Println("incorrect move" + string(bufferBytes))
		conn.Write([]byte(constants.INCORRECT_MOVE_ERROR + "\n"))
		return
	}
	move[1] = strings.Trim(move[1], "\n")
	if len(move[0]) != 2 || len(move[1]) != 2 {
		log.Println("incorrect move" + string(bufferBytes))
		conn.Write([]byte(constants.INCORRECT_MOVE_ERROR + "\n"))
		return
	}

	from, err := strconv.Atoi(move[0])
	if err != nil {
		log.Println(err)
		conn.Write([]byte(constants.INCORRECT_MOVE_ERROR + "\n"))
		return
	}

	to, err := strconv.Atoi(move[1])
	if err != nil {
		log.Println(err)
		conn.Write([]byte(constants.INCORRECT_MOVE_ERROR + "\n"))
		return
	}

	fromX := from / 10
	fromY := from % 10

	toX := to / 10
	toY := to % 10

	err = currRoom.Game.MakeMove(fromX, fromY, toX, toY)
	if err != nil {
		log.Println(err)
		conn.Write([]byte(constants.INCORRECT_MOVE_ERROR + "\n"))
		return
	}
	marshalledGame, err := currRoom.Game.MarshalJSON()
	if err != nil {
		log.Println(err)
		conn.Write([]byte(constants.INCORRECT_MOVE_ERROR + "\n"))
		return
	}
	marshalledGame = append(marshalledGame, []byte("\n")...)
	conn.Write(marshalledGame)
}

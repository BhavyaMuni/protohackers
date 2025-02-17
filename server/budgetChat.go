package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

type BudgetChatServer struct {
	BaseServer
	Users []User
}

type User struct {
	connection *net.Conn
	username   string
}

func NewBudgetChatServer() *BudgetChatServer {
	bcs := &BudgetChatServer{}
	bcs.HandleConnectionFunc = bcs.handleConnection
	bcs.Users = make([]User, 0)
	return bcs
}

func (s *BudgetChatServer) handleConnection(conn net.Conn) {
	log.Println("Connected with...")
	log.Println(conn.RemoteAddr())
	scanner := bufio.NewScanner(conn)
	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
	scanner.Scan()
	name := scanner.Text()
	if !validUsername(name) {
		conn.Write([]byte("Invalid username. Please try again.\n"))
		conn.Close()
		return
	}
	currUsers := s.listUsers()
	s.Users = append(s.Users, User{connection: &conn, username: name})
	conn.Write([]byte("* The room contains: " + strings.Join(currUsers, ", ") + "\n"))
	s.broadcast("* "+name+" has entered the room\n", &conn)

	for scanner.Scan() {
		message := scanner.Text()
		messageToSend := fmt.Sprintf("[%s] %s\n", name, message)
		s.broadcast(messageToSend, &conn)
	}

	defer s.userLeft(&conn)
}

func (s *BudgetChatServer) broadcast(message string, currConn *net.Conn) {
	for _, user := range s.Users {
		if user.connection != currConn {
			(*user.connection).Write([]byte(message))
		}
	}
}

func (s *BudgetChatServer) listUsers() []string {
	var users []string
	for _, user := range s.Users {
		users = append(users, user.username)
	}
	return users
}

func (s *BudgetChatServer) userLeft(conn *net.Conn) {
	for i, user := range s.Users {
		if user.connection == conn {
			s.broadcast("* "+user.username+" has left the room\n", conn)
			s.Users = append(s.Users[:i], s.Users[i+1:]...)
			break
		}
	}
}

func validUsername(name string) bool {
	rx := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return len(name) >= 1 && rx.MatchString(name)
}

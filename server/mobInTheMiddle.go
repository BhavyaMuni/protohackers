package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

type MobInTheMiddleServer struct {
	BaseServer
}

func NewMobInTheMiddleServer() *MobInTheMiddleServer {
	ms := &MobInTheMiddleServer{}
	ms.HandleConnectionFunc = ms.handleConnection
	return ms
}

const UPSTREAM_HOST = "chat.protohackers.com"
const UPSTREAM_PORT = "16963"
const TONY_ADDRESS = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

func (ms MobInTheMiddleServer) handleConnection(conn net.Conn) {
	log.Println("Connected with client:", conn.RemoteAddr())

	// Connect to the upstream server
	upstreamConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", UPSTREAM_HOST, UPSTREAM_PORT))
	if err != nil {
		log.Println("Error connecting to upstream server:", err)
		return
	}
	defer upstreamConn.Close()
	log.Println("Connected to upstream server:", upstreamConn.RemoteAddr())

	// Handle messages from upstream to client
	go ms.handleUpstream(upstreamConn, conn)

	// Handle messages from client to upstream using Scanner
	// Scanner will only return complete lines that end with a newline
	// This ensures we only send complete messages to the upstream server
	reader := bufio.NewReader(conn)
	for {
		// ReadString reads until the first occurrence of the delimiter
		// It will only return when a complete line is available
		message, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from client:", err)
			}
			break
		}

		// Trim the newline character for processing
		message = strings.TrimSuffix(message, "\n")

		// Modify and send the message
		modifiedMessage := modifyMessage(message)
		_, err = upstreamConn.Write([]byte(modifiedMessage + "\n"))
		if err != nil {
			log.Println("Error writing to upstream:", err)
			return
		}
	}
}

func (ms *MobInTheMiddleServer) handleUpstream(upstreamConn net.Conn, downstreamConn net.Conn) {
	upstreamScanner := bufio.NewScanner(upstreamConn)
	for upstreamScanner.Scan() {
		message := upstreamScanner.Text()
		modifiedMessage := modifyMessage(message)
		_, err := downstreamConn.Write([]byte(modifiedMessage + "\n"))
		if err != nil {
			log.Println("Error writing to client:", err)
			return
		}
	}
}

func modifyMessage(message string) string {
	log.Println("Modifying Message:", message)
	// Use lookahead and lookbehind to ensure we only match words surrounded by spaces or string boundaries
	rx := regexp.MustCompile(`^7[a-zA-Z0-9]{25,34}$`)

	// Process the matches - we need to extract the actual match from the submatch array
	var actualMatches []string
	for _, match := range strings.Split(message, " ") {
		if rx.MatchString(match) {
			actualMatches = append(actualMatches, match)
		}
	}
	for _, match := range actualMatches {
		message = strings.Replace(message, match, TONY_ADDRESS, 1)
	}
	return message
}

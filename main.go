package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

var redisEncodings []string = []string{"*", "$"}

// Commands that needs arguments and the number argument needed
var argsCommands map[string]int = map[string]int{"GET": 1, "SET": 2}

func main() {
	store := NewStore()

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn, store)
	}
}

func handleCommand(cmd string, args []string, conn net.Conn, store *Store) error {
	fmt.Println("handling cmd", cmd, args)
	switch cmd {
	case "PING":
		if _, err := conn.Write([]byte(FormatResponse("pong"))); err != nil {
			return err
		}
		return nil
	case "SET":
		// TODO: check args
		store.Set(args[0], args[1])
		for i, arg := range args {
			fmt.Println(i, arg)
		}
		if _, err := conn.Write([]byte(FormatResponse("ok"))); err != nil {
			return err
		}
		return nil
	case "GET":
		v := store.Get(args[0])
		fmt.Println(v)
		if _, err := conn.Write([]byte(FormatResponse(v))); err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

func handleConnection(conn net.Conn, store *Store) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	var lastCmd string
	args := make([]string, 0)
	currArg := 0

	for scanner.Scan() {
		txt := scanner.Text()
		if isRedisEncodingToken(txt) {
			continue
		}
		if nbArgs, exists := argsCommands[strings.ToUpper(txt)]; exists {
			lastCmd = strings.ToUpper(txt)
			args = make([]string, nbArgs)
			continue
		}
		if lastCmd != "" {
			args[currArg] = txt
			currArg++
			if cap(args) != currArg {
				continue
			}
		}
		cmd := strings.ToUpper(txt)
		if lastCmd != "" {
			cmd = lastCmd
			lastCmd = ""
		}

		handleCommand(cmd, args, conn, store)
		args = make([]string, 0)
	}
}

func isRedisEncodingToken(txt string) bool {
	startWith := StartsWith(txt)
	if slices.Contains(redisEncodings, startWith) {
		return true
	}
	return false
}

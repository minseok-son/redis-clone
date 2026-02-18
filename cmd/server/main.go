package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/minseok-son/redis-clone/internal/resp"
	"github.com/minseok-son/redis-clone/internal/storage"
)

func main() {
	listener, _ := net.Listen("tcp", ":6379")
	fmt.Println("Server is listening on port 6379...")
	db := storage.NewDB()
	db.StartJanitor(time.Second)

	for {
		conn, _ := listener.Accept()
		go handleConnection(conn, db)
	}
}

func handleConnection(c net.Conn, db *storage.DB) {
	defer c.Close()
	reader := bufio.NewReader(c)

	for {
		parts, err := resp.Parse(reader)
		if err != nil {
			return
		}

		handleCommand(parts, c, db)
	}
}

func handleCommand(parts []string, c net.Conn, db *storage.DB) {
	if len(parts) == 0 {
		return
	}

	switch strings.ToUpper(parts[0]) {
	case "GET":
		if len(parts) != 2 {
			c.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
			return
		}
		val, ok := db.Get(parts[1])
		if !ok {
			c.Write([]byte("$-1\r\n"))
			return
		}
		c.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
	case "SET":
		if len(parts) == 3 {
			db.Set(parts[1], parts[2], 0)
			c.Write([]byte("+OK\r\n"))
		} else if len(parts) == 5 && strings.ToUpper(parts[3]) == "EX" {
			s, err := strconv.ParseInt(parts[4], 10, 64)
			if err != nil {
				c.Write([]byte("-ERR value is not an integer or out of range\r\n"))
				return
			}

			expiresAt := time.Now().Add(time.Duration(s) * time.Second).UnixNano()

			db.Set(parts[1], parts[2], expiresAt)
			c.Write([]byte("+OK\r\n"))
		} else {
			c.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
			return
		}
	case "DEL":
		if len(parts) != 2 {
			c.Write([]byte("-ERR wrong number of arguments for 'del' command\r\n"))
			return
		}
		
		count := db.Del(parts[1]) // This should return an int (1 if deleted, 0 if not)
		
		// Redis RESP format for integers starts with ':'
		response := fmt.Sprintf(":%d\r\n", count)
		c.Write([]byte(response))
	case "PING":
		c.Write([]byte("+PONG\r\n"))
	default:
		c.Write([]byte("-ERR unknown command\r\n"))
	}
}
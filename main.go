package main

import (
	"io"
	"log"
	"net"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

func Handler(conn *net.TCPConn) {
	defer conn.Close()
	io.WriteString(conn, "Socket Connection!!\n")
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		input := string(buffer[:n])
		switch input {
		case "hello\n":
			_, err = io.WriteString(conn, "こんにちは\n")
		case "bye\n":
			_, err = io.WriteString(conn, "さようなら\n")
			if err == nil {
				return
			}
		default:
			response := "不明なコマンドです:" + input + "\n"
			_, err = io.WriteString(conn, response)
		}
		if err != nil {
			return
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	_db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer _db.Close()

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		err := ln.SetDeadline(time.Now().Add(time.Second * 10))
		if err != nil {
			log.Fatal(err)
		}
		conn, err := ln.AcceptTCP()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			log.Fatal(err)
		}
		go Handler(conn)
	}
}

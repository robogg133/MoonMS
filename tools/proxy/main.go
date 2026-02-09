package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	listenAddr = "127.0.0.1:25565"
	serverAddr = "127.0.0.1:3023"
)

func main() {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Proxy escutando em", listenAddr)

	for {
		clientConn, err := ln.Accept()
		if err != nil {
			log.Println("Erro ao aceitar:", err)
			continue
		}

		go handleClient(clientConn)
	}
}

func handleClient(client net.Conn) {
	defer client.Close()

	server, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Println("error connecting to the server:", err)
		return
	}
	defer server.Close()

	log.Println("Client connected:", client.RemoteAddr())

	// client -> server
	go pipe("C → S", client, server)

	// server -> client
	pipe("S → C", server, client)
}

func pipe(label string, src net.Conn, dst net.Conn) {
	buf := make([]byte, 4096)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			data := buf[:n]

			fmt.Printf("[%s] %d bytes\n%s\n\n",
				label,
				n,
				hex.Dump(data),
			)

			_, _ = dst.Write(data)
		}

		if err != nil {
			if err != io.EOF {
				log.Println("err:", err)
			}
			return
		}
	}
}

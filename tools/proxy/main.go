package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var listenAddr, serverAddr string

func main() {

	flag.StringVar(&listenAddr, "listen-address", "127.0.0.1:25565", "Address that the proxy will listen")
	flag.StringVar(&serverAddr, "server-address", "127.0.0.1:3032", "Address that will be proxied")
	flag.Parse()

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening all connections on", listenAddr)

	for {
		clientConn, err := ln.Accept()
		if err != nil {
			log.Println("error accepting:", err)
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

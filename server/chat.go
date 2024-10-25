package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
)

type client chan<- string // an outgoing message channel
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)
	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func accept(_ context.Context, c chan net.Conn, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		c <- conn
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan net.Conn)
	defer cancel()
	go accept(ctx, c, listener)

	for {
		select {
		case <-ctx.Done():
			return
		case conn := <-c:
			go handleConn(conn)
		}
	}
}

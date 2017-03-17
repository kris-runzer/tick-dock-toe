package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	bind := ":3000"

	game := &Game{}
	game.Reset()

	mux := http.NewServeMux()
	mux.Handle("/", NewLoggingMiddlewareHandler(IndexHandlerFunc))
	mux.Handle("/state", NewStateHandler(game))
	mux.Handle("/move", NewLoggingMiddlewareHandler(NewMakeMoveHandler(game)))
	mux.Handle("/new", NewLoggingMiddlewareHandler(NewNewGameHandler(game)))

	server := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	listener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalln(err)
		os.Exit(2)
	}

	errCh := make(chan error)

	go func() {
		errCh <- server.Serve(listener)
	}()

	log.Printf("[INFO] server started: %s\n", listener.Addr())

	select {
	case signal := <-interrupt():
		log.Printf("[INFO] server interrupted: %s\n", signal)
	case err := <-errCh:
		log.Printf("[ERROR] server failed: %v\n", err)
	}
}

func interrupt() chan os.Signal {
	signalch := make(chan os.Signal)
	signal.Notify(signalch, syscall.SIGINT, syscall.SIGTERM)
	return signalch
}

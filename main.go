package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	bind := flag.String("bind", ":3000", "the http binding port")
	flag.Parse()

	game := &Game{}
	game.Reset()

	mux := http.NewServeMux()

	mw := newLoggingMiddlewareHandlerFunc

	mux.Handle("/", mw(indexHandlerFunc))
	mux.Handle("/state", mw(newStateHandlerFunc(game)))
	mux.Handle("/move", mw(newMakeMoveHandlerFunc(game)))
	mux.Handle("/new", mw(newNewGameHandlerFunc(game)))

	server := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	listener, err := net.Listen("tcp", *bind)
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

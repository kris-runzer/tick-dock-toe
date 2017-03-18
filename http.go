package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kris-runzer/tick-dock-toe/assets"
)

// HTTP Methods
const (
	MethodGet  = "GET"
	MethodPost = "POST"
	MethodPut  = "PUT"
)

var diff = diffFn

func diffFn(start time.Time) time.Duration {
	return time.Now().Sub(start)
}

// NewLoggingMiddlewareHandlerFunc ...
func NewLoggingMiddlewareHandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sw := &statusResponseWriter{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		defer func(start time.Time) {
			log.Printf("[INFO] %s [%d] %s\n", r.RequestURI, sw.Status, diff(start))
		}(time.Now())

		next(sw, r)
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (w *statusResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *statusResponseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

// NewNewGameHandlerFunc ....
func NewNewGameHandlerFunc(game *Game) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		game.Reset()

		responseModel := DefaultResponseModel{
			Board:    game.Board,
			Player:   game.Player,
			NumMoves: game.NumMoves,
			Status:   game.Status,
		}

		if err := json.NewEncoder(w).Encode(responseModel); err != nil {
			jsonErrResponse(w, err)
			return
		}
	}
}

// DefaultResponseModel ...
type DefaultResponseModel struct {
	Board    [3][3]int `json:"board"`
	Player   int       `json:"player"`
	NumMoves int       `json:"numMoves"`
	Status   string    `json:"status"`
}

// NewStateHandlerFunc ...
func NewStateHandlerFunc(game *Game) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		responseModel := DefaultResponseModel{
			Board:    game.Board,
			Player:   game.Player,
			NumMoves: game.NumMoves,
			Status:   game.Status,
		}

		if err := json.NewEncoder(w).Encode(responseModel); err != nil {
			jsonErrResponse(w, err)
			return
		}
	}
}

// NewMakeMoveHandlerFunc ...
func NewMakeMoveHandlerFunc(game *Game) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var model MoveModel

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
			jsonErrResponse(w, err)
			return
		}

		if err := game.MakeMove(model.X, model.Y); err != nil {
			jsonErrResponse(w, err)
			return
		}

		responseModel := DefaultResponseModel{
			Board:    game.Board,
			Player:   game.Player,
			NumMoves: game.NumMoves,
			Status:   game.Status,
		}

		if err := json.NewEncoder(w).Encode(responseModel); err != nil {
			jsonErrResponse(w, err)
			return
		}
	}
}

// MoveModel ...
type MoveModel struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func jsonErrResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	_ = json.NewEncoder(w).Encode(ErrResponseModel{
		Err: err.Error(),
	})
}

// ErrResponseModel ...
type ErrResponseModel struct {
	Err string `json:"error"`
}

// IndexHandlerFunc ...
func IndexHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	html, err := assets.Asset("index.html")
	if err != nil {
		log.Println("[ERROR] failed to retrieve asset index")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	w.Write(html)
}

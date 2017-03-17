package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// HTTP Methods
const (
	MethodGet  = "GET"
	MethodPost = "POST"
	MethodPut  = "PUT"
)

// NewLoggingMiddlewareHandler ...
func NewLoggingMiddlewareHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sw := &statusResponseWriter{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		defer func(start time.Time) {
			log.Printf("[INFO] %s [%d] %s\n", r.RequestURI, sw.Status, time.Now().Sub(start))
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

// NewNewGameHandler ....
func NewNewGameHandler(game *Game) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		game.Reset()

		responseModel := DefaultResponseModel{
			Board:    game.State.Board,
			Player:   game.State.Player,
			NumMoves: game.State.NumMoves,
			Status:   game.State.Status,
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

// NewStateHandler ....
func NewStateHandler(game *Game) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		responseModel := DefaultResponseModel{
			Board:    game.State.Board,
			Player:   game.State.Player,
			NumMoves: game.State.NumMoves,
			Status:   game.State.Status,
		}

		if err := json.NewEncoder(w).Encode(responseModel); err != nil {
			jsonErrResponse(w, err)
			return
		}
	}
}

// NewMakeMoveHandler ...
func NewMakeMoveHandler(game *Game) http.HandlerFunc {
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
			Board:    game.State.Board,
			Player:   game.State.Player,
			NumMoves: game.State.NumMoves,
			Status:   game.State.Status,
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

	html, err := assetsIndexHtmlBytes()
	if err != nil {
		// TODO: error page
		log.Println("[ERROR] failed to retrieve asset index")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	w.Write(html)
}

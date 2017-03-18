package main

import (
	"reflect"
	"testing"
)

func TestPlayerOneWin(t *testing.T) {
	game := &Game{}
	game.Reset()

	_ = game.MakeMove(1, 1)
	_ = game.MakeMove(0, 1)
	_ = game.MakeMove(1, 0)
	_ = game.MakeMove(1, 2)
	_ = game.MakeMove(0, 2)
	_ = game.MakeMove(2, 0)
	_ = game.MakeMove(2, 2)
	_ = game.MakeMove(2, 1)
	_ = game.MakeMove(0, 0)

	expectedState := &State{
		Board:    testNew3x3Board(1, 2, 1, 1, 1, 2, 2, 2, 1),
		Player:   1,
		NumMoves: 9,
		Status:   "end",
	}

	if !reflect.DeepEqual(expectedState, game.State) {
		t.Errorf("unexpected state: %#v", game.State)
	}
}

func TestPlayerTwoWin(t *testing.T) {
	game := &Game{}
	game.Reset()

	_ = game.MakeMove(0, 0)
	_ = game.MakeMove(0, 1)
	_ = game.MakeMove(0, 2)
	_ = game.MakeMove(1, 1)
	_ = game.MakeMove(2, 0)
	_ = game.MakeMove(2, 1)

	expectedState := &State{
		Board:    testNew3x3Board(1, 2, 1, 0, 2, 0, 1, 2, 0),
		Player:   2,
		NumMoves: 6,
		Status:   "end",
	}

	if !reflect.DeepEqual(expectedState, game.State) {
		t.Errorf("unexpected state: %#v", game.State)
	}
}

func TestDraw(t *testing.T) {
	game := &Game{}
	game.Reset()

	_ = game.MakeMove(0, 0)
	_ = game.MakeMove(0, 1)
	_ = game.MakeMove(0, 2)
	_ = game.MakeMove(1, 0)
	_ = game.MakeMove(1, 1)
	_ = game.MakeMove(2, 2)
	_ = game.MakeMove(2, 1)
	_ = game.MakeMove(2, 0)
	_ = game.MakeMove(1, 2)

	expectedState := &State{
		Board:    testNew3x3Board(1, 2, 1, 2, 1, 1, 2, 1, 2),
		Player:   1,
		NumMoves: 9,
		Status:   "draw",
	}

	if !reflect.DeepEqual(expectedState, game.State) {
		t.Errorf("unexpected state: %#v", game.State)
	}
}

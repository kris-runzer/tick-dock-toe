package main

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestGame_Reset(t *testing.T) {
	game := &Game{}
	game.Reset()

	if game == nil {
		t.Fatalf("expected state to be instantiated")
	}

	expectedBoard := [3][3]int{}

	if board := game.Board; !reflect.DeepEqual(expectedBoard, board) {
		t.Errorf("unexpected board: %#v", board)
	}

	if player := game.Player; 1 != player {
		t.Error("unexpected player:", player)
	}

	if numMoves := game.NumMoves; 0 != numMoves {
		t.Error("unexpected numMoves:", numMoves)
	}

	if status := game.Status; "alive" != status {
		t.Error("unexpected status:", status)
	}
}

func TestGame_MakeMove_Noop(t *testing.T) {
	game := &Game{Status: StatusDraw}

	if err := game.MakeMove(0, 0); "game over" != err.Error() {
		t.Error("unexpected err:", err)
	}

	game = &Game{Status: StatusEnd}

	if err := game.MakeMove(0, 0); "game over" != err.Error() {
		t.Error("unexpected err:", err)
	}
}

func TestGame_MakeMove_InvalidMove(t *testing.T) {
	defer func() {
		isValidMove = isValidMoveFn
	}()

	calledBoard := [3][3]int{}
	calledX := 0
	calledY := 0

	isValidMove = func(board [3][3]int, x, y int) error {
		calledBoard = board
		calledX = x
		calledY = y
		return errors.New("boom")
	}

	board := testNew3x3Board(0, 1, 0, 0, 1, 0, 0, 0, 0)
	game := &Game{Board: board}

	if err := game.MakeMove(1, 2); "boom" != errors.Cause(err).Error() {
		t.Error("unexpected err:", err)
	}

	if !reflect.DeepEqual(board, calledBoard) {
		t.Errorf("unexpected board: %#v", calledBoard)
	}

	if 1 != calledX {
		t.Error("unexpected x:", calledX)
	}

	if 2 != calledY {
		t.Error("unexpected y:", calledY)
	}
}

func TestGame_MakeMove_ReturnsOnIsWin(t *testing.T) {
	defer func() {
		isValidMove = isValidMoveFn
		isWin = isWinFn
	}()

	isValidMove = func(board [3][3]int, x, y int) error {
		return nil
	}

	calledBoard := [3][3]int{}
	calledPlayer := 0

	isWin = func(board [3][3]int, player int) bool {
		calledBoard = board
		calledPlayer = player
		return true
	}

	board := testNew3x3Board(0, 1, 0, 0, 1, 0, 0, 0, 0)
	game := &Game{Board: board, Player: 1}

	if err := game.MakeMove(1, 2); nil != err {
		t.Error("unexpected err:", err)
	}

	expectedBoard := testNew3x3Board(0, 1, 0, 0, 1, 1, 0, 0, 0)
	if !reflect.DeepEqual(expectedBoard, calledBoard) {
		t.Errorf("unexpected board: %#v", calledBoard)
	}

	if 1 != calledPlayer {
		t.Error("unexpected player:", calledPlayer)
	}

	if status := game.Status; "end" != status {
		t.Error("unexpected status", status)
	}
}

func TestGame_MakeMove_ReturnsOnIsDraw(t *testing.T) {
	defer func() {
		isValidMove = isValidMoveFn
		isWin = isWinFn
	}()

	isValidMove = func(board [3][3]int, x, y int) error {
		return nil
	}

	isWin = func(board [3][3]int, player int) bool {
		return false
	}

	board := testNew3x3Board(1, 2, 1, 2, 2, 1, 1, 0, 2)
	game := &Game{Board: board, Player: 1, NumMoves: 8}

	if err := game.MakeMove(2, 1); nil != err {
		t.Error("unexpected err:", err)
	}

	expectedBoard := testNew3x3Board(1, 2, 1, 2, 2, 1, 1, 1, 2)
	if gameBoard := game.Board; !reflect.DeepEqual(expectedBoard, gameBoard) {
		t.Errorf("unexpected board: %#v", gameBoard)
	}

	if player := game.Player; 1 != player {
		t.Error("unexpected player:", player)
	}

	if status := game.Status; "draw" != status {
		t.Error("unexpected status", status)
	}
}

func TestGame_MakeMove_SwitchesPlayer(t *testing.T) {
	defer func() {
		isValidMove = isValidMoveFn
		isWin = isWinFn
	}()

	isValidMove = func(board [3][3]int, x, y int) error {
		return nil
	}

	isWin = func(board [3][3]int, player int) bool {
		return false
	}

	board := testEmpty3x3Board()
	game := &Game{Board: board, Player: 1}

	if err := game.MakeMove(2, 1); nil != err {
		t.Error("unexpected err:", err)
	}
	if err := game.MakeMove(1, 1); nil != err {
		t.Error("unexpected err:", err)
	}

	expectedBoard := testNew3x3Board(0, 0, 0, 0, 2, 0, 0, 1, 0)
	if gameBoard := game.Board; !reflect.DeepEqual(expectedBoard, gameBoard) {
		t.Errorf("unexpected board: %#v", gameBoard)
	}

	if player := game.Player; 1 != player {
		t.Error("unexpected player:", player)
	}
}

func testNew3x3Board(x0y0, x0y1, x0y2, x1y0, x1y1, x1y2, x2y0, x2y1, x2y2 int) [3][3]int {
	return [3][3]int{
		[3]int{x0y0, x0y1, x0y2},
		[3]int{x1y0, x1y1, x1y2},
		[3]int{x2y0, x2y1, x2y2},
	}
}

func testEmpty3x3Board() [3][3]int {
	return testNew3x3Board(0, 0, 0, 0, 0, 0, 0, 0, 0)
}

func TestIsValidMove_InvalidXIndex(t *testing.T) {
	board := testEmpty3x3Board()

	if err := isValidMove(board, -1, 0); err.Error() != "invalid x index: -1" {
		t.Error("unexpected err:", err)
	}

	if err := isValidMove(board, 3, 0); err.Error() != "invalid x index: 3" {
		t.Error("unexpected err:", err)
	}
}

func TestIsValidMove_InvalidYIndex(t *testing.T) {
	board := testEmpty3x3Board()

	if err := isValidMove(board, 0, -1); err.Error() != "invalid y index: -1" {
		t.Error("unexpected err:", err)
	}

	if err := isValidMove(board, 0, 3); err.Error() != "invalid y index: 3" {
		t.Error("unexpected err:", err)
	}
}

func TestIsValidMove_NotValidPosition(t *testing.T) {
	tests := []struct {
		InBoard     [3][3]int
		InX         int
		InY         int
		ExpectedErr string
	}{
		{
			InBoard:     testNew3x3Board(1, 0, 0, 0, 0, 0, 0, 0, 0),
			InX:         0,
			InY:         0,
			ExpectedErr: "space already taken: [0][0]: 1",
		},
		{
			InBoard:     testNew3x3Board(0, 1, 0, 0, 0, 0, 0, 0, 0),
			InX:         0,
			InY:         1,
			ExpectedErr: "space already taken: [0][1]: 1",
		},
		{
			InBoard:     testNew3x3Board(0, 0, 1, 0, 0, 0, 0, 0, 0),
			InX:         0,
			InY:         2,
			ExpectedErr: "space already taken: [0][2]: 1",
		},
		{
			InBoard:     testNew3x3Board(0, 0, 0, 1, 0, 0, 0, 0, 0),
			InX:         1,
			InY:         0,
			ExpectedErr: "space already taken: [1][0]: 1",
		},
		{
			InBoard:     testNew3x3Board(0, 0, 0, 0, 1, 0, 0, 0, 0),
			InX:         1,
			InY:         1,
			ExpectedErr: "space already taken: [1][1]: 1",
		},
		{
			InBoard:     testNew3x3Board(0, 0, 0, 0, 0, 1, 0, 0, 0),
			InX:         1,
			InY:         2,
			ExpectedErr: "space already taken: [1][2]: 1",
		},
		{
			InBoard:     testNew3x3Board(0, 0, 0, 0, 0, 0, 1, 0, 0),
			InX:         2,
			InY:         0,
			ExpectedErr: "space already taken: [2][0]: 1",
		},
		{
			InBoard:     testNew3x3Board(0, 0, 0, 0, 0, 0, 0, 1, 0),
			InX:         2,
			InY:         1,
			ExpectedErr: "space already taken: [2][1]: 1",
		},
		{
			InBoard:     testNew3x3Board(0, 0, 0, 0, 0, 0, 0, 0, 1),
			InX:         2,
			InY:         2,
			ExpectedErr: "space already taken: [2][2]: 1",
		},
	}

	for i, test := range tests {
		if err := isValidMove(test.InBoard, test.InX, test.InY); err.Error() != test.ExpectedErr {
			t.Errorf("%d> unexpected err: %v", i, err)
		}
	}
}

func TestIsValidMove_ValidPosition(t *testing.T) {
	tests := []struct {
		InX int
		InY int
	}{
		{InX: 0, InY: 0},
		{InX: 0, InY: 1},
		{InX: 0, InY: 2},
		{InX: 1, InY: 0},
		{InX: 1, InY: 1},
		{InX: 1, InY: 2},
		{InX: 2, InY: 0},
		{InX: 2, InY: 1},
		{InX: 2, InY: 2},
	}

	board := testEmpty3x3Board()

	for i, test := range tests {
		if err := isValidMove(board, test.InX, test.InY); err != nil {
			t.Errorf("%d> unexpected err: %v", i, err)
		}
	}
}

func TestIsWin_NoCheckers(t *testing.T) {
	defer func() {
		winChecks = defaultWinChecks
	}()

	tests := []struct {
		Checks []WinCheck
		OK     bool
	}{
		{
			Checks: nil,
			OK:     false,
		},
		{
			Checks: []WinCheck{
				func([3][3]int, int) bool {
					return false
				},
			},
			OK: false,
		},
		{
			Checks: []WinCheck{
				func([3][3]int, int) bool {
					return true
				},
			},
			OK: true,
		},
		{
			Checks: []WinCheck{
				func([3][3]int, int) bool {
					return false
				},
				func([3][3]int, int) bool {
					return true
				},
			},
			OK: true,
		},
	}

	for i, test := range tests {
		winChecks = test.Checks

		if ok := isWin(testEmpty3x3Board(), 0); ok != test.OK {
			t.Errorf("%d> unexpected win: %t", i, ok)
		}
	}
}

func TestColumnWinCheck(t *testing.T) {
	tests := []struct {
		Board [3][3]int
		IsWin bool
	}{
		{
			Board: testEmpty3x3Board(),
			IsWin: false,
		},
		{
			Board: testNew3x3Board(1, 0, 0, 1, 0, 0, 1, 0, 0),
			IsWin: true,
		},
		{
			Board: testNew3x3Board(0, 1, 0, 0, 1, 0, 0, 1, 0),
			IsWin: true,
		},
		{
			Board: testNew3x3Board(0, 0, 1, 0, 0, 1, 0, 0, 1),
			IsWin: true,
		},
	}

	for i, test := range tests {
		if ok := columnWinCheck(test.Board, 1); ok != test.IsWin {
			t.Errorf("%d> unexpected win: %t", i, ok)
		}
	}
}

func TestRowWinCheck(t *testing.T) {
	tests := []struct {
		Board [3][3]int
		IsWin bool
	}{
		{
			Board: testEmpty3x3Board(),
			IsWin: false,
		},
		{
			Board: testNew3x3Board(1, 1, 1, 0, 0, 0, 0, 0, 0),
			IsWin: true,
		},
		{
			Board: testNew3x3Board(0, 0, 0, 1, 1, 1, 0, 0, 0),
			IsWin: true,
		},
		{
			Board: testNew3x3Board(0, 0, 0, 0, 0, 0, 1, 1, 1),
			IsWin: true,
		},
	}

	for i, test := range tests {
		if ok := rowWinCheck(test.Board, 1); ok != test.IsWin {
			t.Errorf("%d> unexpected win: %t", i, ok)
		}
	}
}

func TestDiagLeftToRightWinCheck(t *testing.T) {
	tests := []struct {
		Board [3][3]int
		IsWin bool
	}{
		{
			Board: testEmpty3x3Board(),
			IsWin: false,
		},
		{
			Board: testNew3x3Board(0, 0, 1, 0, 1, 0, 1, 0, 0),
			IsWin: false,
		},
		{
			Board: testNew3x3Board(1, 0, 0, 0, 1, 0, 0, 0, 1),
			IsWin: true,
		},
	}

	for i, test := range tests {
		if ok := diagLeftToRightWinCheck(test.Board, 1); ok != test.IsWin {
			t.Errorf("%d> unexpected win: %t", i, ok)
		}
	}
}

func TestDiagRightToLeftWinCheck(t *testing.T) {
	tests := []struct {
		Board [3][3]int
		IsWin bool
	}{
		{
			Board: testEmpty3x3Board(),
			IsWin: false,
		},
		{
			Board: testNew3x3Board(1, 0, 0, 0, 1, 0, 0, 0, 1),
			IsWin: false,
		},
		{
			Board: testNew3x3Board(0, 0, 1, 0, 1, 0, 1, 0, 0),
			IsWin: true,
		},
	}

	for i, test := range tests {
		if ok := diagRightToLeftWinCheck(test.Board, 1); ok != test.IsWin {
			t.Errorf("%d> unexpected win: %t", i, ok)
		}
	}
}

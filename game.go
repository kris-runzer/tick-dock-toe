package main

import "github.com/pkg/errors"

// var game *Game

// Game ...
type Game struct {
	State *State
}

// NewGame ...
func NewGame() *Game {
	return &Game{
		State: &State{
			Board:  [3][3]int{},
			Player: 1,
		},
	}
}

// State ...
type State struct {
	Board    [3][3]int
	Player   int
	NumMoves int
}

var (
	// ErrGameOver ...
	ErrGameOver = errors.New("game over")

	// ErrGameDraw ...
	ErrGameDraw = errors.New("game draw")
)

// MakeMove ...
func (g *Game) MakeMove(x, y int) error {
	if err := isValidMoveFn(g.State.Board, x, y); err != nil {
		return errors.Wrap(err, "invalid move")
	}

	g.State.NumMoves++

	if g.State.NumMoves == 9 {
		return ErrGameDraw
	}

	g.State.Board[x][y] = g.State.Player

	if isWin(g.State.Board, g.State.Player) {
		return ErrGameOver
	}

	switch g.State.Player {
	case 1:
		g.State.Player = 2
	case 2:
		g.State.Player = 1
	}

	return nil
}

var isValidMove = isValidMoveFn

func isValidMoveFn(board [3][3]int, x, y int) error {
	if x < 0 || x > 2 {
		return errors.Errorf("invalid x index: %d", x)
	}

	if y < 0 || y > 2 {
		return errors.Errorf("invalid y index: %d", y)
	}

	if val := board[x][y]; val != 0 {
		return errors.Errorf("space already taken: [%d][%d]: %d", x, y, val)
	}

	return nil
}

var isWin = isWinFn

func isWinFn(board [3][3]int, player int) bool {
	for i := 0; i < len(winChecks); i++ {
		if ok := winChecks[i](board, player); ok {
			return true
		}
	}

	return false
}

// WinCheck ...
type WinCheck func(board [3][3]int, player int) bool

var winChecks = defaultWinChecks

var defaultWinChecks = []WinCheck{
	columnWinCheck,
	rowWinCheck,
	diagLeftToRightWinCheck,
	diagRightToLeftWinCheck,
}

// TODO Add comment hear explaining this
var magicNumberWeights = [3][3]int{
	[3]int{8, 1, 6},
	[3]int{3, 5, 7},
	[3]int{4, 9, 2},
}

// columnWinCheck determines if the player has won on a column using
// magic numbers.
func columnWinCheck(board [3][3]int, player int) bool {
	for y := 0; y < 3; y++ {
		sum := 0

		for x := 0; x < 3; x++ {
			if board[x][y] == player {
				sum += magicNumberWeights[x][y]
			}
		}

		if sum == 15 {
			return true
		}
	}

	return false
}

// rowWinCheck determines if the player has won on a row using
// magic numbers.
func rowWinCheck(board [3][3]int, player int) bool {
	for x := 0; x < 3; x++ {
		sum := 0

		for y := 0; y < 3; y++ {
			if board[x][y] == player {
				sum += magicNumberWeights[x][y]
			}
		}

		if sum == 15 {
			return true
		}
	}

	return false
}

// diagLeftToRightWinCheck determines if the player has won on the left-to-right
// diagonal line using magic numbers.
func diagLeftToRightWinCheck(board [3][3]int, player int) bool {
	sum := 0

	for i := 0; i < 3; i++ {
		if board[i][i] == player {
			sum += magicNumberWeights[i][i]
		}
	}

	if sum == 15 {
		return true
	}

	return false
}

// diagRightToLeftWinCheck determines if the player has won on the right-to-left
// diagonal line using magic numbers.
func diagRightToLeftWinCheck(board [3][3]int, player int) bool {
	sum := 0

	for i := 0; i < 3; i++ {
		if board[i][2-i] == player {
			sum += magicNumberWeights[i][2-i]
		}
	}

	if sum == 15 {
		return true
	}

	return false
}

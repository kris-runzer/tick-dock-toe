package main

import "testing"

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

package model

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/rand"
)

// Game struct to hold the current state
type Game struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Board    [][]bool           `json:"-" bson:"board"`     // true for mine, false for safe
	Moves    [][]bool           `json:"moves" bson:"moves"` // true for revealed, false for hidden
	Mines    int                `json:"mines" bson:"mines"`
	Revealed int                `json:"revealed" bson:"revealed"`
	GameOver bool               `json:"game_over" bson:"game_over"`
	Payout   float64            `json:"payout" bson:"payout"`
}

// Updated Payout list derived from the new image
var payoutList = [][]float64{
	{1.01, 1.08, 1.12, 1.18, 1.24, 1.30, 1.37, 1.46, 1.55, 1.65, 1.77, 1.90, 2.06, 2.25, 2.47, 2.75, 3.09, 3.54, 4.12, 4.95, 6.19, 8.25, 12.37, 24.75},
	{1.08, 1.17, 1.29, 1.41, 1.56, 1.74, 1.94, 2.18, 2.47, 2.83, 3.26, 3.81, 4.5, 5.4, 6.6, 8.25, 10.61, 14.14, 19.8, 29.7, 49.5, 99, 297},
	{1.12, 1.29, 1.48, 1.71, 2.00, 2.35, 2.79, 3.35, 4.07, 5.00, 6.26, 7.96, 10.35, 13.8, 18.97, 27.11, 40.66, 65.06, 113.85, 227.7, 569.3, 2277},
	{1.18, 1.41, 1.71, 2.05, 2.58, 3.23, 4.09, 5.26, 6.88, 9.17, 12.51, 17.52, 25.35, 37.95, 59.64, 99.39, 178.91, 357.81, 834.9, 2504, 12523},
	{1.24, 1.56, 2.00, 2.58, 3.39, 4.52, 6.14, 8.5, 12.04, 17.52, 26.77, 40.87, 66.41, 113.85, 208.72, 417.45, 939.26, 2504, 12523},
	{1.3, 1.74, 2.35, 3.23, 4.32, 5.82, 8.12, 11.87, 18.05, 28.1, 45.02, 75.21, 128.89, 232.64, 451.31, 1086.45, 2705.62, 6764.05, 16910.12},
	{1.37, 1.94, 2.79, 4.09, 5.62, 7.74, 10.69, 15.52, 24.31, 37.34, 57.38, 94.72, 153.52, 249.97, 412.92, 709.96, 1320.73, 2458.69, 5107.16, 11363.92},
	{1.46, 2.18, 3.35, 5.26, 8.64, 14.17, 24.47, 44.05, 83.2, 176.8, 356.56, 603.45, 1111.84, 2146.38, 4746.9, 11106.6, 27408.72, 71262.68, 188346.9},
	{1.55, 2.47, 4.07, 6.88, 12.04, 21.89, 41.6, 83.2, 176.8, 404.1, 1010.1, 2828.8, 9193, 36773, 202254, 2022545, 3236072, 4852483, 7291924, 10938483},
	{1.65, 2.83, 5.0, 9.17, 17.52, 33.83, 73.95, 166.4, 404.1, 1010.1, 2828.8, 9193, 36773, 202254, 2022545, 3236072, 4852483, 7291924, 10938483},
	{1.77, 3.26, 6.26, 12.51, 26.27, 58.38, 138.66, 356.56, 1010, 3232, 12123, 56574, 396022, 5148297, 5148297},
	{1.99, 3.81, 7.95, 17.52, 40.87, 102.17, 277.33, 831.98, 2828, 11314, 56574, 396022, 5148297},
	{2.06, 4.5, 10.35, 25.3, 66.41, 189.75, 600.87, 2163, 9193, 49031, 367735, 5148297},
	{2.25, 5.4, 13.8, 37.95, 113.85, 379.5, 1442, 6489, 36773, 294188, 4412826},
	{2.47, 6.6, 18.97, 59.64, 208.72, 834.9, 3965, 23794, 118973, 2022545, 3236072},
	{2.75, 8.25, 27.11, 99.39, 417.45, 2087, 13219, 118973, 2022545},
	{3.09, 10.61, 40.66, 178.91, 939.26, 6261, 59486, 1070759},
	{3.54, 14.14, 65.06, 357.81, 2504, 25047, 475893},
	{4.12, 19.8, 113.9, 834.9, 8766, 175329},
	{4.95, 29.7, 227.7, 2504, 52598},
	{6.19, 49.5, 569.3, 12523},
	{8.25, 99, 297},
	{12.38, 297},
	{24.75},
}

const (
	ROWS    = 5
	COLUMNS = 5
)

func NewGame(mines int) (*Game, error) {
	game := &Game{
		ID:       primitive.NewObjectID(),
		Board:    make([][]bool, ROWS),
		Moves:    make([][]bool, ROWS),
		Mines:    mines,
		Revealed: 0,
		GameOver: false,
		Payout:   1.0,
	}

	// Initialize the board and moves
	for i := range game.Board {
		game.Board[i] = make([]bool, COLUMNS)
		game.Moves[i] = make([]bool, COLUMNS)
	}

	// If the number of mines is greater than the number of cells, return an error (minimum 1 mine)
	if mines > (ROWS*COLUMNS-1) || mines < 1 {
		return nil, errors.New("invalid number of mines")
	}

	// Place mines on the board
	game.placeMines(game, mines)
	return game, nil
}

// PlaceMines randomly places mines on the board
func (g *Game) placeMines(game *Game, mines int) {
	// Seed the math/rand package with a cryptographically secure random number generator
	var b [256]byte
	if _, err := crypto_rand.Read(b[:]); err != nil {
		// TODO: Handle error
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(uint64(binary.LittleEndian.Uint64(b[:])))

	count := 0
	for count < mines {
		row := rand.Intn(ROWS)
		col := rand.Intn(COLUMNS)

		if !game.Board[row][col] {
			game.Board[row][col] = true
			count++
		}
	}
}

// MakeMove reveals a cell on the board
func (g *Game) MakeMove(row, col int) (*Game, error) {
	if g.GameOver {
		return nil, errors.New("game is over. No more moves allowed")

	}

	if row < 0 || row >= ROWS || col < 0 || col >= COLUMNS {
		return nil, errors.New("invalid move. Out of bounds")
	}

	if g.Moves[row][col] {
		return nil, errors.New("cell already revealed")
	}

	g.Moves[row][col] = true

	if g.Board[row][col] {
		g.GameOver = true
		g.Payout = 0
		return g, nil
	} else {
		g.Revealed++
		g.calculatePayout()
		return g, nil
	}
}

// calculatePayout calculates the payout based on the number of mines revealed
func (g *Game) calculatePayout() {
	if g.Revealed > 0 && g.Revealed <= len(payoutList[g.Mines-1]) {
		g.Payout = payoutList[g.Mines-1][g.Revealed-1]
	} else {
		g.Payout = 1.0
	}
}

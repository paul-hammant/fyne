// Package main demonstrates a basic TicTacToe implementation.
// This version works but lacks accessibility features - it cannot be
// navigated via keyboard and provides no feedback for screen readers.
package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Game state
var (
	board       [9]string
	currentTurn = "X"
	gameOver    = false
	cells       [9]*canvas.Text
	cellBgs     [9]*canvas.Rectangle
)

func main() {
	a := app.New()
	w := a.NewWindow("TicTacToe")
	w.Resize(fyne.NewSize(300, 350))

	grid := createGrid(w)
	resetBtn := widget.NewButton("New Game", func() {
		resetGame()
	})

	content := container.NewBorder(nil, resetBtn, nil, nil, grid)
	w.SetContent(content)
	w.ShowAndRun()
}

// createGrid builds the 3x3 game board using simple tappable rectangles
func createGrid(w fyne.Window) *fyne.Container {
	grid := container.NewGridWithColumns(3)

	for i := 0; i < 9; i++ {
		idx := i // Capture for closure

		// Background rectangle
		bg := canvas.NewRectangle(color.RGBA{220, 220, 220, 255})
		bg.SetMinSize(fyne.NewSize(80, 80))
		bg.StrokeColor = color.Black
		bg.StrokeWidth = 2
		cellBgs[idx] = bg

		// Text showing X or O
		text := canvas.NewText("", color.Black)
		text.TextSize = 48
		text.Alignment = fyne.TextAlignCenter
		cells[idx] = text

		// Simple tappable cell - NO keyboard support!
		cell := newTappableCell(idx, w, bg, text)
		grid.Add(cell)
	}

	return grid
}

// tappableCell is a simple tappable container - mouse only, no accessibility
type tappableCell struct {
	widget.BaseWidget
	idx    int
	window fyne.Window
	bg     *canvas.Rectangle
	text   *canvas.Text
}

func newTappableCell(idx int, w fyne.Window, bg *canvas.Rectangle, text *canvas.Text) *tappableCell {
	cell := &tappableCell{idx: idx, window: w, bg: bg, text: text}
	cell.ExtendBaseWidget(cell)
	return cell
}

func (c *tappableCell) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewStack(c.bg, container.NewCenter(c.text)),
	)
}

// Tapped handles mouse clicks only
func (c *tappableCell) Tapped(*fyne.PointEvent) {
	if gameOver || board[c.idx] != "" {
		return
	}

	board[c.idx] = currentTurn
	c.text.Text = currentTurn
	c.text.Color = getPlayerColor(currentTurn)
	c.text.Refresh()

	if winner := checkWinner(); winner != "" {
		gameOver = true
		dialog.ShowInformation("Game Over", winner+" wins!", c.window)
		return
	}

	if isBoardFull() {
		gameOver = true
		dialog.ShowInformation("Game Over", "It's a draw!", c.window)
		return
	}

	// Switch turns
	if currentTurn == "X" {
		currentTurn = "O"
	} else {
		currentTurn = "X"
	}
}

func getPlayerColor(player string) color.Color {
	if player == "X" {
		return color.RGBA{0, 0, 200, 255} // Blue
	}
	return color.RGBA{200, 0, 0, 255} // Red
}

func checkWinner() string {
	// Winning combinations
	lines := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}

	for _, line := range lines {
		if board[line[0]] != "" &&
			board[line[0]] == board[line[1]] &&
			board[line[1]] == board[line[2]] {
			return board[line[0]]
		}
	}
	return ""
}

func isBoardFull() bool {
	for _, cell := range board {
		if cell == "" {
			return false
		}
	}
	return true
}

func resetGame() {
	board = [9]string{}
	currentTurn = "X"
	gameOver = false

	for _, text := range cells {
		text.Text = ""
		text.Refresh()
	}
}

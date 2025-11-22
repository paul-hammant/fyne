// Package main demonstrates an accessible TicTacToe implementation.
// This version includes full keyboard navigation, focus indicators,
// status announcements, and follows accessibility best practices.
package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Ensure our cell implements the required interfaces
var (
	_ fyne.Focusable     = (*accessibleCell)(nil)
	_ fyne.Tappable      = (*accessibleCell)(nil)
	_ desktop.Hoverable  = (*accessibleCell)(nil)
	_ desktop.Keyable    = (*accessibleCell)(nil)
)

// Game holds all game state in a clean, testable structure
type Game struct {
	board       [9]string
	currentTurn string
	gameOver    bool
	cells       [9]*accessibleCell
	statusLabel *widget.Label
	window      fyne.Window
}

// NewGame creates a new game instance
func NewGame(w fyne.Window) *Game {
	return &Game{
		currentTurn: "X",
		window:      w,
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("TicTacToe (Accessible)")
	w.Resize(fyne.NewSize(340, 420))

	game := NewGame(w)
	content := game.buildUI()
	w.SetContent(content)

	// Set initial focus to center cell for better UX
	w.Canvas().Focus(game.cells[4])

	w.ShowAndRun()
}

// buildUI creates the complete game interface
func (g *Game) buildUI() fyne.CanvasObject {
	// Status label announces game state - critical for screen readers
	g.statusLabel = widget.NewLabel("Player X's turn. Use Tab to navigate, Space or Enter to place.")
	g.statusLabel.Wrapping = fyne.TextWrapWord
	g.statusLabel.Alignment = fyne.TextAlignCenter

	// Build the accessible grid
	grid := g.createAccessibleGrid()

	// Reset button with clear labeling
	resetBtn := widget.NewButton("New Game (R)", func() {
		g.resetGame()
	})

	// Instructions for accessibility
	instructions := widget.NewLabel("Keyboard: Tab/Arrows to move, Space/Enter to play")
	instructions.Alignment = fyne.TextAlignCenter

	return container.NewBorder(
		container.NewVBox(g.statusLabel, widget.NewSeparator()),
		container.NewVBox(widget.NewSeparator(), instructions, resetBtn),
		nil, nil,
		grid,
	)
}

// createAccessibleGrid builds the 3x3 game board with full accessibility
func (g *Game) createAccessibleGrid() *fyne.Container {
	grid := container.NewGridWithColumns(3)

	for i := 0; i < 9; i++ {
		cell := newAccessibleCell(i, g)
		g.cells[i] = cell
		grid.Add(cell)
	}

	return grid
}

// accessibleCell is a fully accessible game cell
type accessibleCell struct {
	widget.BaseWidget

	idx     int
	game    *Game
	value   string
	focused bool
	hovered bool

	// Visual elements
	bg         *canvas.Rectangle
	focusRing  *canvas.Rectangle
	text       *canvas.Text
}

func newAccessibleCell(idx int, game *Game) *accessibleCell {
	cell := &accessibleCell{
		idx:  idx,
		game: game,
	}
	cell.ExtendBaseWidget(cell)
	return cell
}

// CreateRenderer sets up the visual representation
func (c *accessibleCell) CreateRenderer() fyne.WidgetRenderer {
	// Focus ring - visible when focused for keyboard users
	c.focusRing = canvas.NewRectangle(color.Transparent)
	c.focusRing.StrokeColor = color.RGBA{0, 120, 215, 255} // Accessible blue
	c.focusRing.StrokeWidth = 3

	// Background
	c.bg = canvas.NewRectangle(color.RGBA{240, 240, 240, 255})
	c.bg.StrokeColor = color.RGBA{100, 100, 100, 255}
	c.bg.StrokeWidth = 1

	// Text for X/O
	c.text = canvas.NewText("", color.Black)
	c.text.TextSize = 48
	c.text.Alignment = fyne.TextAlignCenter

	content := container.NewStack(
		c.focusRing,
		container.NewPadded(c.bg),
		container.NewCenter(c.text),
	)

	return &accessibleCellRenderer{
		cell:    c,
		content: content,
	}
}

// MinSize ensures cells are large enough for touch and visibility
func (c *accessibleCell) MinSize() fyne.Size {
	return fyne.NewSize(90, 90) // WCAG recommends 44x44 minimum touch target
}

// --- Focusable Interface Implementation ---

// FocusGained is called when this cell receives keyboard focus
func (c *accessibleCell) FocusGained() {
	c.focused = true
	c.Refresh()
}

// FocusLost is called when this cell loses keyboard focus
func (c *accessibleCell) FocusLost() {
	c.focused = false
	c.Refresh()
}

// TypedRune handles character input (not used for TicTacToe)
func (c *accessibleCell) TypedRune(r rune) {
	// 'r' for reset from any cell
	if r == 'r' || r == 'R' {
		c.game.resetGame()
	}
}

// TypedKey handles special key presses for navigation and activation
func (c *accessibleCell) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeySpace, fyne.KeyReturn, fyne.KeyEnter:
		// Activate the cell (same as clicking)
		c.activate()

	case fyne.KeyUp:
		c.moveFocus(-3) // Move up one row
	case fyne.KeyDown:
		c.moveFocus(3) // Move down one row
	case fyne.KeyLeft:
		c.moveFocus(-1) // Move left
	case fyne.KeyRight:
		c.moveFocus(1) // Move right

	case fyne.KeyHome:
		// Jump to first cell
		c.game.window.Canvas().Focus(c.game.cells[0])
	case fyne.KeyEnd:
		// Jump to last cell
		c.game.window.Canvas().Focus(c.game.cells[8])
	}
}

// moveFocus moves focus to an adjacent cell with bounds checking
func (c *accessibleCell) moveFocus(delta int) {
	newIdx := c.idx + delta

	// Bounds checking
	if newIdx < 0 || newIdx > 8 {
		return
	}

	// Prevent wrapping at row edges for left/right movement
	if delta == -1 && c.idx%3 == 0 {
		return // Already at left edge
	}
	if delta == 1 && c.idx%3 == 2 {
		return // Already at right edge
	}

	c.game.window.Canvas().Focus(c.game.cells[newIdx])
}

// --- Desktop Keyable Interface (for physical keyboard support) ---

func (c *accessibleCell) KeyDown(ev *fyne.KeyEvent) {
	// Handled by TypedKey
}

func (c *accessibleCell) KeyUp(ev *fyne.KeyEvent) {
	// Not needed
}

// --- Tappable Interface Implementation ---

// Tapped handles mouse/touch input
func (c *accessibleCell) Tapped(*fyne.PointEvent) {
	// Also grab focus when tapped - important for mixed input users
	c.game.window.Canvas().Focus(c)
	c.activate()
}

// --- Hoverable Interface Implementation ---

// MouseIn provides visual feedback on hover
func (c *accessibleCell) MouseIn(*desktop.MouseEvent) {
	c.hovered = true
	c.Refresh()
}

// MouseMoved is required by the interface
func (c *accessibleCell) MouseMoved(*desktop.MouseEvent) {}

// MouseOut removes hover visual feedback
func (c *accessibleCell) MouseOut() {
	c.hovered = false
	c.Refresh()
}

// --- Game Logic ---

// activate places the current player's mark
func (c *accessibleCell) activate() {
	if c.game.gameOver || c.value != "" {
		// Provide feedback for invalid moves
		if c.value != "" {
			c.game.updateStatus(fmt.Sprintf("Cell already taken by %s. Choose another.", c.value))
		}
		return
	}

	// Place the mark
	c.value = c.game.currentTurn
	c.game.board[c.idx] = c.value
	c.Refresh()

	// Check for winner
	if winner := c.game.checkWinner(); winner != "" {
		c.game.gameOver = true
		c.game.updateStatus(fmt.Sprintf("Game Over! Player %s wins! Press R for new game.", winner))
		c.game.highlightWinningLine(winner)
		return
	}

	// Check for draw
	if c.game.isBoardFull() {
		c.game.gameOver = true
		c.game.updateStatus("Game Over! It's a draw! Press R for new game.")
		return
	}

	// Switch turns
	if c.game.currentTurn == "X" {
		c.game.currentTurn = "O"
	} else {
		c.game.currentTurn = "X"
	}

	c.game.updateStatus(fmt.Sprintf("Player %s's turn.", c.game.currentTurn))
}

// getPositionName returns a human-readable position for screen readers
func (c *accessibleCell) getPositionName() string {
	positions := []string{
		"top-left", "top-center", "top-right",
		"middle-left", "center", "middle-right",
		"bottom-left", "bottom-center", "bottom-right",
	}
	return positions[c.idx]
}

// updateStatus updates the status label (announces to screen readers)
func (g *Game) updateStatus(msg string) {
	g.statusLabel.SetText(msg)
}

func (g *Game) checkWinner() string {
	lines := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}

	for _, line := range lines {
		if g.board[line[0]] != "" &&
			g.board[line[0]] == g.board[line[1]] &&
			g.board[line[1]] == g.board[line[2]] {
			return g.board[line[0]]
		}
	}
	return ""
}

func (g *Game) highlightWinningLine(winner string) {
	lines := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	for _, line := range lines {
		if g.board[line[0]] == winner &&
			g.board[line[1]] == winner &&
			g.board[line[2]] == winner {
			// Highlight winning cells
			for _, idx := range line {
				g.cells[idx].bg.FillColor = color.RGBA{144, 238, 144, 255} // Light green
				g.cells[idx].Refresh()
			}
			return
		}
	}
}

func (g *Game) isBoardFull() bool {
	for _, cell := range g.board {
		if cell == "" {
			return false
		}
	}
	return true
}

func (g *Game) resetGame() {
	g.board = [9]string{}
	g.currentTurn = "X"
	g.gameOver = false

	for _, cell := range g.cells {
		cell.value = ""
		cell.bg.FillColor = color.RGBA{240, 240, 240, 255}
		cell.Refresh()
	}

	g.updateStatus("New game! Player X's turn. Use Tab to navigate, Space or Enter to place.")

	// Return focus to center cell
	g.window.Canvas().Focus(g.cells[4])
}

// --- Custom Renderer ---

type accessibleCellRenderer struct {
	cell    *accessibleCell
	content *fyne.Container
}

func (r *accessibleCellRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
}

func (r *accessibleCellRenderer) MinSize() fyne.Size {
	return r.cell.MinSize()
}

func (r *accessibleCellRenderer) Refresh() {
	// Update text
	r.cell.text.Text = r.cell.value
	if r.cell.value == "X" {
		r.cell.text.Color = color.RGBA{0, 0, 180, 255} // Blue - good contrast
	} else if r.cell.value == "O" {
		r.cell.text.Color = color.RGBA{180, 0, 0, 255} // Red - good contrast
	}

	// Update focus ring visibility
	if r.cell.focused {
		r.cell.focusRing.StrokeColor = color.RGBA{0, 120, 215, 255}
		r.cell.focusRing.StrokeWidth = 3
	} else {
		r.cell.focusRing.StrokeColor = color.Transparent
		r.cell.focusRing.StrokeWidth = 0
	}

	// Update hover state (subtle background change)
	if r.cell.hovered && !r.cell.game.gameOver && r.cell.value == "" {
		r.cell.bg.FillColor = color.RGBA{220, 235, 250, 255} // Light blue hint
	} else if r.cell.bg.FillColor != (color.RGBA{144, 238, 144, 255}) { // Don't override winning highlight
		r.cell.bg.FillColor = color.RGBA{240, 240, 240, 255}
	}

	r.cell.focusRing.Refresh()
	r.cell.bg.Refresh()
	r.cell.text.Refresh()
	canvas.Refresh(r.cell)
}

func (r *accessibleCellRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

func (r *accessibleCellRenderer) Destroy() {}

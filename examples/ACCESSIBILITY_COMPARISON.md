# TicTacToe Accessibility Comparison

A developer education guide comparing `tictactoe.go` (basic) vs `tictactoe-accessible.go` (accessible).

## Executive Summary

| Aspect | Basic Version | Accessible Version |
|--------|--------------|-------------------|
| Mouse/Touch | Yes | Yes |
| Keyboard Navigation | No | Yes (Tab, Arrows) |
| Keyboard Activation | No | Yes (Space, Enter) |
| Focus Indicators | No | Yes (visible ring) |
| Status Announcements | Dialog only | Continuous label |
| Touch Target Size | 80x80 | 90x90 (WCAG AA) |
| Code Complexity | ~150 lines | ~350 lines |

---

## 1. Interface Implementation

### Basic Version: Minimal Interfaces

```go
// Only implements Tappable - mouse/touch only
type tappableCell struct {
    widget.BaseWidget
    // ...
}

func (c *tappableCell) Tapped(*fyne.PointEvent) {
    // Handle click
}
```

**Problem**: Users who can't use a mouse (motor impairments, screen reader users) cannot play the game.

### Accessible Version: Complete Interface Set

```go
// Compile-time interface verification
var (
    _ fyne.Focusable     = (*accessibleCell)(nil)
    _ fyne.Tappable      = (*accessibleCell)(nil)
    _ desktop.Hoverable  = (*accessibleCell)(nil)
    _ desktop.Keyable    = (*accessibleCell)(nil)
)
```

**Key Interfaces**:

| Interface | Purpose | Methods |
|-----------|---------|---------|
| `fyne.Focusable` | Keyboard focus | `FocusGained()`, `FocusLost()`, `TypedKey()`, `TypedRune()` |
| `fyne.Tappable` | Mouse/touch | `Tapped()` |
| `desktop.Hoverable` | Visual hover feedback | `MouseIn()`, `MouseOut()`, `MouseMoved()` |
| `desktop.Keyable` | Raw key events | `KeyDown()`, `KeyUp()` |

---

## 2. Keyboard Navigation

### Basic Version: None

```go
// No keyboard support - users must use mouse
```

### Accessible Version: Full Arrow Key + Tab Support

```go
func (c *accessibleCell) TypedKey(ev *fyne.KeyEvent) {
    switch ev.Name {
    case fyne.KeySpace, fyne.KeyReturn, fyne.KeyEnter:
        c.activate()  // Place mark

    case fyne.KeyUp:
        c.moveFocus(-3)  // Up one row
    case fyne.KeyDown:
        c.moveFocus(3)   // Down one row
    case fyne.KeyLeft:
        c.moveFocus(-1)  // Left one cell
    case fyne.KeyRight:
        c.moveFocus(1)   // Right one cell

    case fyne.KeyHome:
        c.game.window.Canvas().Focus(c.game.cells[0])
    case fyne.KeyEnd:
        c.game.window.Canvas().Focus(c.game.cells[8])
    }
}
```

**Navigation Features**:
- **Tab**: Standard focus navigation (built into Fyne)
- **Arrow Keys**: Grid-aware directional movement
- **Home/End**: Jump to first/last cell
- **Space/Enter**: Activate (place mark)
- **R key**: Reset game from any cell

**Bounds Checking** prevents impossible moves:

```go
func (c *accessibleCell) moveFocus(delta int) {
    newIdx := c.idx + delta
    if newIdx < 0 || newIdx > 8 {
        return  // Out of grid
    }
    // Prevent horizontal wrapping
    if delta == -1 && c.idx%3 == 0 {
        return  // At left edge
    }
    if delta == 1 && c.idx%3 == 2 {
        return  // At right edge
    }
    c.game.window.Canvas().Focus(c.game.cells[newIdx])
}
```

---

## 3. Visual Focus Indicators

### Basic Version: No Indication

Users cannot see which cell has focus (if any).

### Accessible Version: Clear Focus Ring

```go
func (r *accessibleCellRenderer) Refresh() {
    if r.cell.focused {
        r.cell.focusRing.StrokeColor = color.RGBA{0, 120, 215, 255}
        r.cell.focusRing.StrokeWidth = 3
    } else {
        r.cell.focusRing.StrokeColor = color.Transparent
        r.cell.focusRing.StrokeWidth = 0
    }
}
```

**WCAG 2.4.7**: Focus must be visible. The 3px blue ring provides:
- High contrast against the gray background
- Sufficient thickness to be noticeable
- Consistent appearance across all cells

---

## 4. Status Communication

### Basic Version: Modal Dialogs Only

```go
if winner := checkWinner(); winner != "" {
    gameOver = true
    dialog.ShowInformation("Game Over", winner+" wins!", c.window)
}
```

**Problems**:
- No turn indication until game ends
- Dialogs interrupt flow
- Screen readers may miss dialog content

### Accessible Version: Live Status Region

```go
// Persistent status label
g.statusLabel = widget.NewLabel("Player X's turn. Use Tab to navigate...")

// Updated throughout game
func (g *Game) updateStatus(msg string) {
    g.statusLabel.SetText(msg)  // Announces to assistive tech
}

// Examples of status updates:
"Player X's turn."
"Cell already taken by O. Choose another."
"Game Over! Player X wins! Press R for new game."
```

**Benefits**:
- Continuous game state awareness
- Instructions always visible
- Non-blocking feedback
- Works with screen readers (live region pattern)

---

## 5. Touch Target Sizing

### Basic Version: 80x80 pixels

```go
bg.SetMinSize(fyne.NewSize(80, 80))
```

### Accessible Version: 90x90 pixels

```go
func (c *accessibleCell) MinSize() fyne.Size {
    return fyne.NewSize(90, 90)  // WCAG recommends 44x44 minimum
}
```

**WCAG 2.5.5 (AAA)**: Target size should be at least 44x44 CSS pixels. The accessible version exceeds this, improving usability for:
- Users with motor impairments
- Touch screen users
- Users with tremors

---

## 6. Mixed Input Support

### Basic Version: Separate Input Modes

Mouse and keyboard are disconnected experiences.

### Accessible Version: Unified Experience

```go
func (c *accessibleCell) Tapped(*fyne.PointEvent) {
    // Grab focus when tapped - bridges mouse and keyboard
    c.game.window.Canvas().Focus(c)
    c.activate()
}
```

Users can:
1. Click a cell to play AND focus it
2. Continue with keyboard from that position
3. Switch between input methods freely

---

## 7. Hover States

### Basic Version: None

No visual feedback on hover.

### Accessible Version: Subtle Hover Indication

```go
func (c *accessibleCell) MouseIn(*desktop.MouseEvent) {
    c.hovered = true
    c.Refresh()
}

// In renderer:
if r.cell.hovered && !r.cell.game.gameOver && r.cell.value == "" {
    r.cell.bg.FillColor = color.RGBA{220, 235, 250, 255}  // Light blue
}
```

**Benefits**:
- Indicates interactive elements
- Shows valid move targets
- Improves discoverability

---

## 8. Error Feedback

### Basic Version: Silent Failure

```go
func (c *tappableCell) Tapped(*fyne.PointEvent) {
    if gameOver || board[c.idx] != "" {
        return  // Silently ignored
    }
}
```

### Accessible Version: Explicit Feedback

```go
func (c *accessibleCell) activate() {
    if c.value != "" {
        c.game.updateStatus(fmt.Sprintf(
            "Cell already taken by %s. Choose another.", c.value))
        return
    }
}
```

Users know WHY their action failed.

---

## 9. Code Organization

### Basic Version: Global State

```go
var (
    board       [9]string
    currentTurn = "X"
    gameOver    = false
    cells       [9]*canvas.Text
)
```

**Problems**: Hard to test, maintain, or extend.

### Accessible Version: Encapsulated Game State

```go
type Game struct {
    board       [9]string
    currentTurn string
    gameOver    bool
    cells       [9]*accessibleCell
    statusLabel *widget.Label
    window      fyne.Window
}

func NewGame(w fyne.Window) *Game {
    return &Game{currentTurn: "X", window: w}
}
```

**Benefits**:
- Testable (inject mock window)
- Multiple games possible
- Clear ownership of state
- Easier to add features

---

## 10. Key Takeaways for Developers

### Always Implement

1. **`fyne.Focusable`** - Required for keyboard users
2. **Visual focus indicators** - Required by WCAG 2.4.7
3. **Keyboard activation** - Space/Enter should work like click

### Consider Implementing

4. **`desktop.Hoverable`** - Improves mouse UX
5. **Arrow key navigation** - For grid/list widgets
6. **Status announcements** - For dynamic content

### Best Practices

7. **Minimum touch targets**: 44x44 pixels (WCAG 2.5.5)
8. **Color contrast**: 4.5:1 for text (WCAG 1.4.3)
9. **Don't rely on color alone**: Use shapes/text too
10. **Test with keyboard only**: Unplug your mouse!

---

## Testing Checklist

- [ ] Can complete entire game using only keyboard?
- [ ] Is focused cell always visible?
- [ ] Are all status changes announced?
- [ ] Do colors have sufficient contrast?
- [ ] Are touch targets large enough?
- [ ] Does hover provide feedback?
- [ ] Are error states communicated?

---

## Resources

- [Fyne Accessibility Interfaces](https://docs.fyne.io/)
- [WCAG 2.1 Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [Go Widget Development](https://docs.fyne.io/extend/custom-widget)

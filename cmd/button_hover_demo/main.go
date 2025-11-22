package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// InteractiveButton is a custom widget that implements all desktop interaction interfaces:
// - desktop.Hoverable: Mouse enter/move/exit events
// - desktop.Mouseable: Mouse button down/up events
// - desktop.Cursorable: Custom cursor appearance
// - desktop.Keyable: Keyboard key down/up events (requires focus)
type InteractiveButton struct {
	widget.Button
}

// NewInteractiveButton creates a new interactive button
func NewInteractiveButton(text string, tapped func()) *InteractiveButton {
	btn := &InteractiveButton{}
	btn.Text = text
	btn.OnTapped = tapped
	btn.ExtendBaseWidget(btn)
	return btn
}

// ========== desktop.Hoverable interface ==========

// MouseIn is called when the mouse pointer enters the button
func (b *InteractiveButton) MouseIn(e *desktop.MouseEvent) {
	fmt.Printf("Hoverable.MouseIn: Mouse entered at (%.2f, %.2f)\n",
		e.Position.X, e.Position.Y)
}

// MouseMoved is called when the mouse pointer moves over the button
func (b *InteractiveButton) MouseMoved(e *desktop.MouseEvent) {
	fmt.Printf("Hoverable.MouseMoved: Mouse moved to (%.2f, %.2f)\n",
		e.Position.X, e.Position.Y)
}

// MouseOut is called when the mouse pointer leaves the button
func (b *InteractiveButton) MouseOut() {
	fmt.Println("Hoverable.MouseOut: Mouse left the button")
}

// ========== desktop.Mouseable interface ==========

// MouseDown is called when a mouse button is pressed on the button
func (b *InteractiveButton) MouseDown(e *desktop.MouseEvent) {
	buttonName := getButtonName(e.Button)
	modifiers := getModifierString(e.Modifier)
	fmt.Printf("Mouseable.MouseDown: %s button pressed at (%.2f, %.2f)%s\n",
		buttonName, e.Position.X, e.Position.Y, modifiers)
}

// MouseUp is called when a mouse button is released on the button
func (b *InteractiveButton) MouseUp(e *desktop.MouseEvent) {
	buttonName := getButtonName(e.Button)
	modifiers := getModifierString(e.Modifier)
	fmt.Printf("Mouseable.MouseUp: %s button released at (%.2f, %.2f)%s\n",
		buttonName, e.Position.X, e.Position.Y, modifiers)
}

// ========== desktop.Cursorable interface ==========

// Cursor returns the cursor type to display when hovering over the button
func (b *InteractiveButton) Cursor() desktop.Cursor {
	// Return a pointer cursor (hand) to indicate the button is clickable
	return desktop.PointerCursor
}

// ========== desktop.Keyable interface (requires fyne.Focusable) ==========

// KeyDown is called when a key is pressed while the button has focus
func (b *InteractiveButton) KeyDown(e *fyne.KeyEvent) {
	modifiers := getModifierString(e.Modifier)
	fmt.Printf("Keyable.KeyDown: Key '%s' pressed%s\n", e.Name, modifiers)
}

// KeyUp is called when a key is released while the button has focus
func (b *InteractiveButton) KeyUp(e *fyne.KeyEvent) {
	modifiers := getModifierString(e.Modifier)
	fmt.Printf("Keyable.KeyUp: Key '%s' released%s\n", e.Name, modifiers)
}

// FocusGained is called when the button gains keyboard focus
func (b *InteractiveButton) FocusGained() {
	fmt.Println("Focusable.FocusGained: Button now has keyboard focus (try pressing keys!)")
}

// FocusLost is called when the button loses keyboard focus
func (b *InteractiveButton) FocusLost() {
	fmt.Println("Focusable.FocusLost: Button lost keyboard focus")
}

// ========== Helper functions ==========

func getButtonName(btn desktop.MouseButton) string {
	switch btn {
	case desktop.MouseButtonPrimary:
		return "Primary (Left)"
	case desktop.MouseButtonSecondary:
		return "Secondary (Right)"
	case desktop.MouseButtonTertiary:
		return "Tertiary (Middle)"
	default:
		return fmt.Sprintf("Button %d", btn)
	}
}

func getModifierString(mod fyne.KeyModifier) string {
	if mod == 0 {
		return ""
	}
	var modifiers []string
	if mod&fyne.KeyModifierShift != 0 {
		modifiers = append(modifiers, "Shift")
	}
	if mod&fyne.KeyModifierControl != 0 {
		modifiers = append(modifiers, "Ctrl")
	}
	if mod&fyne.KeyModifierAlt != 0 {
		modifiers = append(modifiers, "Alt")
	}
	if mod&fyne.KeyModifierSuper != 0 {
		modifiers = append(modifiers, "Super")
	}
	if len(modifiers) > 0 {
		return fmt.Sprintf(" with %v", modifiers)
	}
	return ""
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Desktop Interfaces Demo")
	myWindow.Resize(fyne.NewSize(600, 500))

	// Create an interactive button that implements all desktop interfaces
	interactiveButton := NewInteractiveButton("Interactive Button - Try Everything!", func() {
		fmt.Println("✓ Button.Tapped: Button was clicked!")
	})

	// Create other widgets that do NOT implement desktop interfaces
	regularButton := widget.NewButton("Regular Button (No desktop events)", func() {
		fmt.Println("Regular button clicked")
	})

	label := widget.NewLabel("Regular label (No desktop events)")

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Regular entry (No desktop events)")

	checkbox := widget.NewCheck("Regular checkbox (No desktop events)", func(checked bool) {
		fmt.Printf("Checkbox changed: %v\n", checked)
	})

	separator1 := widget.NewSeparator()
	separator2 := widget.NewSeparator()

	instructions := widget.NewLabel(
		"The top button implements ALL desktop interfaces:\n\n" +
			"• Hoverable: Move mouse in/over/out of button\n" +
			"• Mouseable: Click with left/right/middle buttons\n" +
			"• Cursorable: Shows pointer cursor when hovering\n" +
			"• Keyable: Click to focus, then press keyboard keys\n\n" +
			"Try Shift/Ctrl/Alt with mouse clicks or key presses!\n" +
			"Watch the console for detailed event reporting.\n\n" +
			"The widgets below do NOT implement these interfaces.")
	instructions.Wrapping = fyne.TextWrapWord

	// Layout all widgets
	content := container.NewVBox(
		instructions,
		separator1,
		widget.NewLabel("Interactive button (all desktop interfaces):"),
		interactiveButton,
		separator2,
		widget.NewLabel("Regular widgets (no desktop interfaces):"),
		regularButton,
		label,
		entry,
		checkbox,
	)

	myWindow.SetContent(content)

	fmt.Println("=== Desktop Interfaces Demo Started ===")
	fmt.Println()
	fmt.Println("This demo shows desktop.Hoverable, desktop.Mouseable,")
	fmt.Println("desktop.Cursorable, and desktop.Keyable interfaces.")
	fmt.Println()
	fmt.Println("Try these interactions with the interactive button:")
	fmt.Println("  1. Move mouse in/over/out (Hoverable)")
	fmt.Println("  2. Click with different mouse buttons (Mouseable)")
	fmt.Println("  3. Notice the cursor changes to a pointer (Cursorable)")
	fmt.Println("  4. Click to focus, then press keys (Keyable)")
	fmt.Println("  5. Try modifier keys (Shift/Ctrl/Alt) with clicks or keys")
	fmt.Println()
	fmt.Println("Watch below for event reports:")
	fmt.Println("=====================================")
	fmt.Println()

	myWindow.ShowAndRun()
}

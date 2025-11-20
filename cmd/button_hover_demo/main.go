package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// HoverableButton is a custom widget that wraps a button and implements desktop.Hoverable
type HoverableButton struct {
	widget.Button
}

// NewHoverableButton creates a new hoverable button
func NewHoverableButton(text string, tapped func()) *HoverableButton {
	btn := &HoverableButton{}
	btn.Text = text
	btn.OnTapped = tapped
	btn.ExtendBaseWidget(btn)
	return btn
}

// MouseIn is called when the mouse pointer enters the button
func (h *HoverableButton) MouseIn(e *desktop.MouseEvent) {
	fmt.Printf("MouseIn: Mouse entered the button at position (%.2f, %.2f)\n",
		e.Position.X, e.Position.Y)
}

// MouseMoved is called when the mouse pointer moves over the button
func (h *HoverableButton) MouseMoved(e *desktop.MouseEvent) {
	fmt.Printf("MouseMoved: Mouse moved over the button to position (%.2f, %.2f)\n",
		e.Position.X, e.Position.Y)
}

// MouseOut is called when the mouse pointer leaves the button
func (h *HoverableButton) MouseOut() {
	fmt.Println("MouseOut: Mouse left the button")
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Button Hover Demo")
	myWindow.Resize(fyne.NewSize(500, 400))

	// Create a hoverable button that reports mouse movements
	hoverButton := NewHoverableButton("Hover Over Me!", func() {
		fmt.Println("Button was clicked!")
	})

	// Create other widgets that are NOT hoverable - they won't report mouse movements
	regularButton := widget.NewButton("Regular Button (No hover events)", func() {
		fmt.Println("Regular button clicked")
	})

	label := widget.NewLabel("This is a regular label (No hover events)")

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Regular entry (No hover events)")

	checkbox := widget.NewCheck("Regular checkbox (No hover events)", func(checked bool) {
		fmt.Printf("Checkbox changed: %v\n", checked)
	})

	separator := widget.NewSeparator()

	instructions := widget.NewLabel("Instructions:\n" +
		"1. Move your mouse INTO the 'Hover Over Me!' button\n" +
		"2. Move your mouse AROUND inside the button\n" +
		"3. Move your mouse OUT OF the button\n" +
		"4. Watch the console for hover events!\n\n" +
		"The other widgets do NOT implement desktop.Hoverable,\n" +
		"so they won't report hover events to the console.")
	instructions.Wrapping = fyne.TextWrapWord

	// Layout all widgets
	content := container.NewVBox(
		instructions,
		separator,
		widget.NewLabel("Hoverable button:"),
		hoverButton,
		separator,
		widget.NewLabel("Non-hoverable widgets:"),
		regularButton,
		label,
		entry,
		checkbox,
	)

	myWindow.SetContent(content)

	fmt.Println("=== Button Hover Demo Started ===")
	fmt.Println("Move your mouse over the 'Hover Over Me!' button to see hover events")
	fmt.Println("The other widgets will NOT generate hover events")
	fmt.Println()

	myWindow.ShowAndRun()
}

package tutorials

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// RotatingKnobScreen demonstrates the rotating knob widget with various configurations
func RotatingKnobScreen(_ fyne.Window) fyne.CanvasObject {
	// Basic knob example
	basicKnob := widget.NewRotatingKnob(0, 100)
	basicKnob.SetValue(50)
	basicValueLabel := widget.NewLabel("Value: 50.0")
	basicKnob.OnChanged = func(value float64) {
		basicValueLabel.SetText(fmt.Sprintf("Value: %.1f", value))
	}

	basicCard := widget.NewCard("Basic Knob", "Standard 0-100 range",
		container.NewVBox(
			container.NewCenter(basicKnob),
			basicValueLabel,
		))

	// Temperature control knob
	tempKnob := widget.NewRotatingKnob(-20, 40)
	tempKnob.SetValue(20)
	tempKnob.Step = 0.5
	tempKnob.TickCount = 13 // -20, -15, -10, ..., 35, 40
	tempValueLabel := widget.NewLabel("Temperature: 20.0°C")
	tempKnob.OnChanged = func(value float64) {
		tempValueLabel.SetText(fmt.Sprintf("Temperature: %.1f°C", value))
	}

	tempCard := widget.NewCard("Temperature Control", "Range: -20°C to 40°C",
		container.NewVBox(
			container.NewCenter(tempKnob),
			tempValueLabel,
		))

	// Volume knob (with custom angles)
	volumeKnob := widget.NewRotatingKnob(0, 10)
	volumeKnob.SetValue(5)
	volumeKnob.StartAngle = -135
	volumeKnob.EndAngle = 135
	volumeKnob.Step = 0.1
	volumeKnob.TickCount = 11
	volumeValueLabel := widget.NewLabel("Volume: 5.0")
	volumeKnob.OnChanged = func(value float64) {
		volumeValueLabel.SetText(fmt.Sprintf("Volume: %.1f", value))
	}

	volumeCard := widget.NewCard("Volume Control", "Range: 0-10, 270° sweep",
		container.NewVBox(
			container.NewCenter(volumeKnob),
			volumeValueLabel,
		))

	// Angle selector (wrapping enabled)
	angleKnob := widget.NewRotatingKnob(0, 359)
	angleKnob.SetValue(0)
	angleKnob.Wrapping = true
	angleKnob.StartAngle = 0
	angleKnob.EndAngle = 359
	angleKnob.TickCount = 8 // N, NE, E, SE, S, SW, W, NW
	angleValueLabel := widget.NewLabel("Angle: 0°")
	angleKnob.OnChanged = func(value float64) {
		angleValueLabel.SetText(fmt.Sprintf("Angle: %.0f°", value))
	}

	angleCard := widget.NewCard("Angle Selector", "Wrapping enabled, full circle",
		container.NewVBox(
			container.NewCenter(angleKnob),
			angleValueLabel,
		))

	// Data binding example
	boundData := binding.NewFloat()
	boundData.Set(25.0)
	boundKnob := widget.NewRotatingKnobWithData(0, 100, boundData)
	boundValueLabel := widget.NewLabel("Value: 25.0")
	boundData.AddListener(binding.NewDataListener(func() {
		val, _ := boundData.Get()
		boundValueLabel.SetText(fmt.Sprintf("Value: %.1f", val))
	}))

	// External control buttons for data binding demo
	incButton := widget.NewButton("Increment", func() {
		val, _ := boundData.Get()
		boundData.Set(val + 5)
	})
	decButton := widget.NewButton("Decrement", func() {
		val, _ := boundData.Get()
		boundData.Set(val - 5)
	})

	boundCard := widget.NewCard("Data Binding", "Bound to external data",
		container.NewVBox(
			container.NewCenter(boundKnob),
			boundValueLabel,
			container.NewGridWithColumns(2, incButton, decButton),
		))

	// Disabled knob example
	disabledKnob := widget.NewRotatingKnob(0, 100)
	disabledKnob.SetValue(75)
	disabledKnob.Disable()
	disabledLabel := widget.NewLabel("Value: 75.0 (Disabled)")

	enableToggle := widget.NewCheck("Enable Knob", func(checked bool) {
		if checked {
			disabledKnob.Enable()
			disabledLabel.SetText(fmt.Sprintf("Value: %.1f (Enabled)", disabledKnob.Value))
		} else {
			disabledKnob.Disable()
			disabledLabel.SetText(fmt.Sprintf("Value: %.1f (Disabled)", disabledKnob.Value))
		}
	})

	disabledCard := widget.NewCard("Disabled State", "Toggle to enable/disable",
		container.NewVBox(
			container.NewCenter(disabledKnob),
			disabledLabel,
			enableToggle,
		))

	// Fine control knob (no ticks, small steps)
	fineKnob := widget.NewRotatingKnob(0, 1)
	fineKnob.SetValue(0.5)
	fineKnob.Step = 0.01
	fineKnob.ShowTicks = false
	fineValueLabel := widget.NewLabel("Value: 0.500")
	fineKnob.OnChanged = func(value float64) {
		fineValueLabel.SetText(fmt.Sprintf("Value: %.3f", value))
	}

	fineCard := widget.NewCard("Fine Control", "0-1 range, no ticks, 0.01 step",
		container.NewVBox(
			container.NewCenter(fineKnob),
			fineValueLabel,
		))

	// Interactive test panel
	testKnob := widget.NewRotatingKnob(0, 100)
	testKnob.SetValue(50)
	testValueLabel := widget.NewLabel("Value: 50.0")
	testEventLog := widget.NewLabel("Events: None")

	changedCount := 0
	endedCount := 0

	testKnob.OnChanged = func(value float64) {
		changedCount++
		testValueLabel.SetText(fmt.Sprintf("Value: %.1f", value))
		testEventLog.SetText(fmt.Sprintf("OnChanged: %d | OnChangeEnded: %d", changedCount, endedCount))
	}
	testKnob.OnChangeEnded = func(value float64) {
		endedCount++
		testEventLog.SetText(fmt.Sprintf("OnChanged: %d | OnChangeEnded: %d", changedCount, endedCount))
	}

	setMinButton := widget.NewButton("Set Min (0)", func() {
		testKnob.SetValue(testKnob.Min)
	})
	setMaxButton := widget.NewButton("Set Max (100)", func() {
		testKnob.SetValue(testKnob.Max)
	})
	setMidButton := widget.NewButton("Set Mid (50)", func() {
		testKnob.SetValue((testKnob.Min + testKnob.Max) / 2)
	})
	resetEventsButton := widget.NewButton("Reset Events", func() {
		changedCount = 0
		endedCount = 0
		testEventLog.SetText("Events: Reset")
	})

	testCard := widget.NewCard("Interactive Test", "Test automation features",
		container.NewVBox(
			container.NewCenter(testKnob),
			testValueLabel,
			testEventLog,
			container.NewGridWithColumns(2,
				setMinButton, setMaxButton,
			),
			container.NewGridWithColumns(2,
				setMidButton, resetEventsButton,
			),
		))

	// Instructions
	instructionsText := canvas.NewText(
		"How to use:\n"+
			"• Drag: Click and drag in a circular motion\n"+
			"• Tap: Click at any position to jump to that value\n"+
			"• Keyboard: Use arrow keys (Up/Right increase, Down/Left decrease)\n"+
			"• Keyboard: Home/End for min/max, Page Up/Down for larger steps\n"+
			"• Scroll: Use mouse wheel to adjust value\n"+
			"• Hover: Visual feedback when mouse is over knob\n"+
			"• Focus: Click to focus, visual indicator when focused",
		theme.ForegroundColor(),
	)
	instructionsText.TextSize = 12

	instructions := widget.NewCard("Instructions", "Interaction methods",
		container.NewVBox(instructionsText),
	)

	// Features panel
	featuresText := canvas.NewText(
		"Features:\n"+
			"• Configurable range (min, max)\n"+
			"• Custom start/end angles (partial rotation)\n"+
			"• Optional value wrapping\n"+
			"• Visual tick marks (configurable count)\n"+
			"• Step size for incremental changes\n"+
			"• Hover and focus states\n"+
			"• Disable/enable support\n"+
			"• Data binding support\n"+
			"• OnChanged and OnChangeEnded callbacks\n"+
			"• Full keyboard accessibility\n"+
			"• Touch-friendly interaction",
		theme.ForegroundColor(),
	)
	featuresText.TextSize = 12

	features := widget.NewCard("Features", "Capabilities",
		container.NewVBox(featuresText),
	)

	// Automation example code
	codeText := `// Automated Testing Example
func TestRotatingKnobAutomation(t *testing.T) {
    knob := widget.NewRotatingKnob(0, 100)

    // Set value programmatically
    knob.SetValue(75)
    assert.Equal(t, 75.0, knob.Value)

    // Test callbacks
    var lastValue float64
    knob.OnChanged = func(value float64) {
        lastValue = value
    }

    // Simulate user interactions
    test.Tap(knob)
    test.Drag(knob, startPos, delta)
    test.Type(knob, "Up")

    // Verify state
    assert.True(t, knob.Focused())
    assert.False(t, knob.Disabled())
}`

	codeLabel := widget.NewLabel(codeText)
	automationCard := widget.NewCard("Automation Example", "Testing code",
		container.NewScroll(codeLabel),
	)

	// Layout everything in a grid
	leftColumn := container.NewVBox(
		basicCard,
		tempCard,
		volumeCard,
		angleCard,
	)

	rightColumn := container.NewVBox(
		boundCard,
		disabledCard,
		fineCard,
		testCard,
	)

	bottomRow := container.NewVBox(
		instructions,
		features,
		automationCard,
	)

	topRow := container.NewGridWithColumns(2,
		leftColumn,
		rightColumn,
	)

	return container.NewBorder(
		nil, nil, nil, nil,
		container.NewVScroll(
			container.NewVBox(
				topRow,
				bottomRow,
			),
		),
	)
}

// RotatingKnobTitle returns the title for the rotating knob tutorial
func RotatingKnobTitle() string {
	return "Rotating Knob"
}

// RotatingKnobDescription returns the description for the rotating knob tutorial
func RotatingKnobDescription() string {
	return "Circular dial/knob control for value selection, similar to a potentiometer or volume knob"
}

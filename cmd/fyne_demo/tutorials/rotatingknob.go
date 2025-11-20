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

// makeKnobWithIcon creates a knob with an icon and value display
func makeKnobWithIcon(knob *widget.RotatingKnob, icon string, valueLabel *widget.Label, accentColor color.Color) fyne.CanvasObject {
	// Apply custom color if provided
	if accentColor != nil {
		knob.AccentColor = accentColor
	}

	// Create icon text
	iconText := canvas.NewText(icon, theme.ForegroundColor())
	iconText.TextSize = 32
	iconText.Alignment = fyne.TextAlignCenter

	// Style the value label
	valueLabel.TextStyle = fyne.TextStyle{Bold: true}
	valueLabel.Alignment = fyne.TextAlignCenter

	// Create a visual container with icon above, knob in center, value below
	return container.NewVBox(
		container.NewCenter(iconText),
		layout.NewSpacer(),
		container.NewCenter(knob),
		layout.NewSpacer(),
		container.NewCenter(valueLabel),
	)
}

// RotatingKnobScreen demonstrates the rotating knob widget with various configurations
func RotatingKnobScreen(_ fyne.Window) fyne.CanvasObject {
	// 1. BASIC KNOB - Simple percentage control
	basicKnob := widget.NewRotatingKnob(0, 100)
	basicKnob.SetValue(50)
	basicKnob.AccentColor = color.NRGBA{R: 100, G: 149, B: 237, A: 255} // Cornflower blue
	basicValueLabel := widget.NewLabel("50%")
	basicValueLabel.TextStyle = fyne.TextStyle{Bold: true}
	basicValueLabel.Alignment = fyne.TextAlignCenter

	basicKnob.OnChanged = func(value float64) {
		basicValueLabel.SetText(fmt.Sprintf("%.0f%%", value))
	}

	basicDisplay := makeKnobWithIcon(basicKnob, "📊", basicValueLabel, basicKnob.AccentColor)

	basicCard := widget.NewCard("Basic Knob", "Standard 0-100 range",
		container.NewCenter(basicDisplay))

	// 2. TEMPERATURE CONTROL - Blue/Red gradient feel
	tempKnob := widget.NewRotatingKnob(-20, 40)
	tempKnob.SetValue(20)
	tempKnob.Step = 0.5
	tempKnob.TickCount = 13
	// Use blue-to-red color based on temperature
	tempKnob.AccentColor = color.NRGBA{R: 255, G: 69, B: 0, A: 255} // Red-Orange for warmth
	tempKnob.TrackColor = color.NRGBA{R: 70, G: 130, B: 180, A: 80}  // Steel blue (faded)

	tempValueLabel := widget.NewLabel("20.0°C")
	tempValueLabel.TextStyle = fyne.TextStyle{Bold: true}
	tempValueLabel.Alignment = fyne.TextAlignCenter

	tempKnob.OnChanged = func(value float64) {
		tempValueLabel.SetText(fmt.Sprintf("%.1f°C", value))
		// Change color based on temperature
		if value < 0 {
			tempKnob.AccentColor = color.NRGBA{R: 0, G: 149, B: 255, A: 255} // Cold blue
		} else if value < 20 {
			tempKnob.AccentColor = color.NRGBA{R: 100, G: 200, B: 255, A: 255} // Cool blue
		} else if value < 30 {
			tempKnob.AccentColor = color.NRGBA{R: 255, G: 165, B: 0, A: 255} // Orange
		} else {
			tempKnob.AccentColor = color.NRGBA{R: 255, G: 69, B: 0, A: 255} // Hot red
		}
		tempKnob.Refresh()
	}

	tempDisplay := makeKnobWithIcon(tempKnob, "🌡️", tempValueLabel, nil)

	tempCard := widget.NewCard("Temperature Control", "Dynamic color (-20°C to 40°C)",
		container.NewCenter(tempDisplay))

	// 3. VOLUME CONTROL - Goes to 11! (Spinal Tap reference)
	volumeKnob := widget.NewRotatingKnob(0, 11)
	volumeKnob.SetValue(5)
	volumeKnob.StartAngle = -90  // 270° (left/9 o'clock)
	volumeKnob.EndAngle = 90     // 90° (right/3 o'clock) - 180° sweep
	volumeKnob.Step = 0.5
	volumeKnob.TickCount = 12 // 0-11
	volumeKnob.AccentColor = color.NRGBA{R: 50, G: 205, B: 50, A: 255}    // Lime green
	volumeKnob.WedgeColor = color.NRGBA{R: 50, G: 205, B: 50, A: 60}      // Semi-transparent green wedge
	volumeKnob.TrackColor = color.NRGBA{R: 80, G: 80, B: 80, A: 40}       // Subtle gray track
	volumeKnob.ShowTicks = true

	volumeValueLabel := widget.NewLabel("5")
	volumeValueLabel.TextStyle = fyne.TextStyle{Bold: true, Italic: false}
	volumeValueLabel.Alignment = fyne.TextAlignCenter

	// Create a special "11" indicator
	volume11Label := widget.NewLabel("")
	volume11Label.TextStyle = fyne.TextStyle{Bold: true}
	volume11Label.Alignment = fyne.TextAlignCenter

	volumeKnob.OnChanged = func(value float64) {
		if value == 11 {
			volumeValueLabel.SetText("11")
			volume11Label.SetText("🎸 IT GOES TO ELEVEN! 🎸")
			volumeKnob.AccentColor = color.NRGBA{R: 255, G: 215, B: 0, A: 255}    // Gold for 11!
			volumeKnob.WedgeColor = color.NRGBA{R: 255, G: 215, B: 0, A: 80}      // Gold wedge
		} else {
			volumeValueLabel.SetText(fmt.Sprintf("%.0f", value))
			volume11Label.SetText("")
			// Green intensity increases with volume
			intensity := uint8(50 + (value/11)*205)
			volumeKnob.AccentColor = color.NRGBA{R: 0, G: intensity, B: 0, A: 255}
			volumeKnob.WedgeColor = color.NRGBA{R: 0, G: intensity, B: 0, A: 60} // Matching wedge
		}
		volumeKnob.Refresh()
	}

	volumeIcon := canvas.NewText("🔊", theme.ForegroundColor())
	volumeIcon.TextSize = 32
	volumeIcon.Alignment = fyne.TextAlignCenter

	volumeDisplay := container.NewVBox(
		container.NewCenter(volumeIcon),
		layout.NewSpacer(),
		container.NewCenter(volumeKnob),
		layout.NewSpacer(),
		container.NewCenter(volumeValueLabel),
		container.NewCenter(volume11Label),
	)

	volumeCard := widget.NewCard("Volume Control", "These go to eleven! 🎸",
		container.NewCenter(volumeDisplay))

	// 4. ANGLE SELECTOR - Compass style
	angleKnob := widget.NewRotatingKnob(0, 359)
	angleKnob.SetValue(0)
	angleKnob.Wrapping = true
	angleKnob.StartAngle = 0
	angleKnob.EndAngle = 359
	angleKnob.TickCount = 8 // N, NE, E, SE, S, SW, W, NW
	angleKnob.AccentColor = color.NRGBA{R: 138, G: 43, B: 226, A: 255} // Blue-violet

	angleValueLabel := widget.NewLabel("0° N")
	angleValueLabel.TextStyle = fyne.TextStyle{Bold: true}
	angleValueLabel.Alignment = fyne.TextAlignCenter

	angleKnob.OnChanged = func(value float64) {
		direction := getCompassDirection(value)
		angleValueLabel.SetText(fmt.Sprintf("%.0f° %s", value, direction))
	}

	angleDisplay := makeKnobWithIcon(angleKnob, "🧭", angleValueLabel, angleKnob.AccentColor)

	angleCard := widget.NewCard("Angle Selector", "Full 360° with wrapping",
		container.NewCenter(angleDisplay))

	// 5. DATA BINDING - Purple with sync icon
	boundData := binding.NewFloat()
	boundData.Set(25.0)
	boundKnob := widget.NewRotatingKnobWithData(0, 100, boundData)
	boundKnob.AccentColor = color.NRGBA{R: 147, G: 112, B: 219, A: 255} // Medium purple

	boundValueLabel := widget.NewLabel("25")
	boundValueLabel.TextStyle = fyne.TextStyle{Bold: true}
	boundValueLabel.Alignment = fyne.TextAlignCenter

	boundData.AddListener(binding.NewDataListener(func() {
		val, _ := boundData.Get()
		boundValueLabel.SetText(fmt.Sprintf("%.0f", val))
	}))

	incButton := widget.NewButton("+ 5", func() {
		val, _ := boundData.Get()
		boundData.Set(val + 5)
	})
	decButton := widget.NewButton("- 5", func() {
		val, _ := boundData.Get()
		boundData.Set(val - 5)
	})

	boundIcon := canvas.NewText("🔄", theme.ForegroundColor())
	boundIcon.TextSize = 32
	boundIcon.Alignment = fyne.TextAlignCenter

	boundDisplay := container.NewVBox(
		container.NewCenter(boundIcon),
		layout.NewSpacer(),
		container.NewCenter(boundKnob),
		layout.NewSpacer(),
		container.NewCenter(boundValueLabel),
		container.NewGridWithColumns(2, incButton, decButton),
	)

	boundCard := widget.NewCard("Data Binding", "Bound to external data",
		container.NewCenter(boundDisplay))

	// 6. DISABLED STATE - Gray with lock icon
	disabledKnob := widget.NewRotatingKnob(0, 100)
	disabledKnob.SetValue(75)
	disabledKnob.Disable()

	disabledLabel := widget.NewLabel("75 (Locked)")
	disabledLabel.Alignment = fyne.TextAlignCenter
	disabledLabel.TextStyle = fyne.TextStyle{Bold: true}

	enableToggle := widget.NewCheck("Unlock", func(checked bool) {
		if checked {
			disabledKnob.Enable()
			disabledKnob.AccentColor = color.NRGBA{R: 34, G: 139, B: 34, A: 255} // Forest green
			disabledLabel.SetText(fmt.Sprintf("%.0f (Unlocked)", disabledKnob.Value))
		} else {
			disabledKnob.Disable()
			disabledKnob.AccentColor = nil
			disabledLabel.SetText(fmt.Sprintf("%.0f (Locked)", disabledKnob.Value))
		}
		disabledKnob.Refresh()
	})

	disabledKnob.OnChanged = func(value float64) {
		if !disabledKnob.Disabled() {
			disabledLabel.SetText(fmt.Sprintf("%.0f (Unlocked)", value))
		}
	}

	disabledDisplay := makeKnobWithIcon(disabledKnob, "🔒", disabledLabel, nil)

	disabledCard := widget.NewCard("Disabled State", "Toggle to unlock",
		container.NewVBox(
			container.NewCenter(disabledDisplay),
			container.NewCenter(enableToggle),
		))

	// 7. FINE CONTROL - Cyan, no ticks, precision dial
	fineKnob := widget.NewRotatingKnob(0, 1)
	fineKnob.SetValue(0.5)
	fineKnob.Step = 0.001
	fineKnob.ShowTicks = false
	fineKnob.AccentColor = color.NRGBA{R: 0, G: 206, B: 209, A: 255} // Dark turquoise

	fineValueLabel := widget.NewLabel("0.500")
	fineValueLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	fineValueLabel.Alignment = fyne.TextAlignCenter

	fineKnob.OnChanged = func(value float64) {
		fineValueLabel.SetText(fmt.Sprintf("%.3f", value))
	}

	fineDisplay := makeKnobWithIcon(fineKnob, "🎯", fineValueLabel, fineKnob.AccentColor)

	fineCard := widget.NewCard("Fine Control", "0-1 range, 0.001 step precision",
		container.NewCenter(fineDisplay))

	// 8. INTERACTIVE TEST PANEL - Rainbow colors
	testKnob := widget.NewRotatingKnob(0, 100)
	testKnob.SetValue(50)

	testValueLabel := widget.NewLabel("50")
	testValueLabel.TextStyle = fyne.TextStyle{Bold: true}
	testValueLabel.Alignment = fyne.TextAlignCenter

	testEventLog := widget.NewLabel("Events: None")
	testEventLog.Alignment = fyne.TextAlignCenter

	changedCount := 0
	endedCount := 0

	testKnob.OnChanged = func(value float64) {
		changedCount++
		testValueLabel.SetText(fmt.Sprintf("%.0f", value))
		testEventLog.SetText(fmt.Sprintf("OnChanged: %d | OnChangeEnded: %d", changedCount, endedCount))

		// Rainbow effect based on value
		hue := value / 100.0
		testKnob.AccentColor = hueToRGB(hue)
		testKnob.Refresh()
	}

	testKnob.OnChangeEnded = func(value float64) {
		endedCount++
		testEventLog.SetText(fmt.Sprintf("OnChanged: %d | OnChangeEnded: %d", changedCount, endedCount))
	}

	setMinButton := widget.NewButton("Min (0)", func() {
		testKnob.SetValue(testKnob.Min)
	})
	setMaxButton := widget.NewButton("Max (100)", func() {
		testKnob.SetValue(testKnob.Max)
	})
	setMidButton := widget.NewButton("Mid (50)", func() {
		testKnob.SetValue((testKnob.Min + testKnob.Max) / 2)
	})
	resetEventsButton := widget.NewButton("Reset Events", func() {
		changedCount = 0
		endedCount = 0
		testEventLog.SetText("Events: Reset")
	})

	testIcon := canvas.NewText("🌈", theme.ForegroundColor())
	testIcon.TextSize = 32
	testIcon.Alignment = fyne.TextAlignCenter

	testDisplay := container.NewVBox(
		container.NewCenter(testIcon),
		layout.NewSpacer(),
		container.NewCenter(testKnob),
		layout.NewSpacer(),
		container.NewCenter(testValueLabel),
		testEventLog,
		container.NewGridWithColumns(2, setMinButton, setMaxButton),
		container.NewGridWithColumns(2, setMidButton, resetEventsButton),
	)

	testCard := widget.NewCard("Interactive Test", "Rainbow colors, event tracking",
		container.NewCenter(testDisplay))

	// Instructions with visual styling
	instructionsText := canvas.NewText(
		"✨ INTERACTION GUIDE ✨\n\n"+
			"🖱️  Drag: Click and rotate in circular motion\n"+
			"👆 Tap: Click anywhere to jump to that value\n"+
			"⌨️  Keyboard: Arrow keys (↑/→ increase, ↓/← decrease)\n"+
			"🏠 Home/End: Jump to min/max values\n"+
			"📄 Page Up/Down: Larger steps (10x)\n"+
			"🖲️  Scroll: Mouse wheel for fine adjustment\n"+
			"👀 Hover: Visual feedback when mouse is over\n"+
			"🎯 Focus: Visual indicator when focused",
		theme.ForegroundColor(),
	)
	instructionsText.TextSize = 11
	instructionsText.Alignment = fyne.TextAlignLeading

	instructions := widget.NewCard("How to Use", "Multiple interaction methods",
		container.NewVBox(instructionsText))

	// Features panel with icons
	featuresText := canvas.NewText(
		"🎨 FEATURES:\n\n"+
			"• Custom colors and visual styling\n"+
			"• Configurable value ranges (min, max)\n"+
			"• Custom start/end angles (partial rotation)\n"+
			"• Optional value wrapping (360° controls)\n"+
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
	featuresText.TextSize = 11

	features := widget.NewCard("Widget Capabilities", "Production-ready features",
		container.NewVBox(featuresText))

	// Layout in a responsive grid
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

// getCompassDirection returns the compass direction for a given angle
func getCompassDirection(angle float64) string {
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((angle + 22.5) / 45) % 8
	return directions[index]
}

// hueToRGB converts a hue value (0-1) to RGB color
func hueToRGB(hue float64) color.Color {
	// Simple HSV to RGB with S=1, V=1
	h := hue * 6.0
	x := uint8(255 * (1 - abs(mod(h, 2.0)-1)))

	switch int(h) {
	case 0:
		return color.NRGBA{R: 255, G: x, B: 0, A: 255}
	case 1:
		return color.NRGBA{R: x, G: 255, B: 0, A: 255}
	case 2:
		return color.NRGBA{R: 0, G: 255, B: x, A: 255}
	case 3:
		return color.NRGBA{R: 0, G: x, B: 255, A: 255}
	case 4:
		return color.NRGBA{R: x, G: 0, B: 255, A: 255}
	default:
		return color.NRGBA{R: 255, G: 0, B: x, A: 255}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func mod(x, y float64) float64 {
	return x - y*float64(int(x/y))
}

// RotatingKnobTitle returns the title for the rotating knob tutorial
func RotatingKnobTitle() string {
	return "Rotating Knob"
}

// RotatingKnobDescription returns the description for the rotating knob tutorial
func RotatingKnobDescription() string {
	return "Circular dial/knob control for value selection with rich visual customization"
}

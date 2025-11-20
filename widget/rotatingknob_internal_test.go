package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestRotatingKnobRenderer_Objects(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)

	objects := renderer.Objects()

	// Should have at least: track, active, indicator, thumb, centerDot
	assert.GreaterOrEqual(t, len(objects), 5)

	// Check object types
	assert.IsType(t, &canvas.Arc{}, objects[0])    // track
	assert.IsType(t, &canvas.Arc{}, objects[1])    // active
	assert.IsType(t, &canvas.Line{}, objects[2])   // indicator
	assert.IsType(t, &canvas.Circle{}, objects[3]) // thumb
	assert.IsType(t, &canvas.Circle{}, objects[4]) // centerDot
}

func TestRotatingKnobRenderer_ObjectsWithTicks(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	knob.ShowTicks = true
	knob.TickCount = 5

	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)
	objects := renderer.Objects()

	// Should have 5 base objects + 5 tick marks
	assert.Equal(t, 10, len(objects))

	// Last 5 objects should be tick lines
	for i := 5; i < 10; i++ {
		assert.IsType(t, &canvas.Line{}, objects[i])
	}
}

func TestRotatingKnobRenderer_ObjectsWithoutTicks(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	knob.ShowTicks = false

	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)
	objects := renderer.Objects()

	// Should have only 5 base objects
	assert.Equal(t, 5, len(objects))
}

func TestRotatingKnobRenderer_Layout(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)

	size := fyne.NewSize(100, 100)
	renderer.Layout(size)

	// Track should be slightly smaller
	assert.Less(t, renderer.track.Size().Width, float32(100))
	assert.Less(t, renderer.track.Size().Height, float32(100))

	// Active should be same size as track
	assert.Equal(t, renderer.track.Size(), renderer.active.Size())
}

func TestRotatingKnobRenderer_MinSize(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)

	minSize := renderer.MinSize()

	// Minimum size should be at least 60x60
	assert.GreaterOrEqual(t, minSize.Width, float32(60))
	assert.GreaterOrEqual(t, minSize.Height, float32(60))
}

func TestRotatingKnobRenderer_Refresh(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)

	// Initial state
	renderer.Refresh()
	assert.NotNil(t, renderer.indicator.StrokeColor)
}

func TestRotatingKnobRenderer_RefreshDisabled(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	knob.Disable()

	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)
	renderer.Refresh()

	// When disabled, colors should be disabled theme color
	assert.Equal(t, theme.DisabledColor(), renderer.indicator.StrokeColor)
	assert.Equal(t, theme.DisabledColor(), renderer.thumb.FillColor)
}

func TestRotatingKnobRenderer_RefreshHovered(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	knob.hovered = true

	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)
	renderer.Refresh()

	// When hovered, active color should be hover color
	assert.Equal(t, theme.HoverColor(), renderer.active.StrokeColor)
	assert.Equal(t, theme.HoverColor(), renderer.thumb.FillColor)
}

func TestRotatingKnobRenderer_RefreshFocused(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	knob.focused = true

	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)
	renderer.Refresh()

	// When focused, indicator should be thicker
	assert.Equal(t, float32(4), renderer.indicator.StrokeWidth)
}

func TestRotatingKnobRenderer_IndicatorPosition(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	knob.SetValue(0) // Minimum value

	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)
	size := fyne.NewSize(100, 100)
	renderer.Layout(size)

	// Indicator should start from center
	centerX := float32(50)
	centerY := float32(50)
	assert.Equal(t, centerX, renderer.indicator.Position1.X)
	assert.Equal(t, centerY, renderer.indicator.Position1.Y)

	// Position2 should be different (pointing somewhere)
	assert.NotEqual(t, centerX, renderer.indicator.Position2.X)
}

func TestRotatingKnobRenderer_IndicatorRotation(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)
	size := fyne.NewSize(100, 100)

	// Set to min value
	knob.SetValue(0)
	renderer.Layout(size)
	minPos := renderer.indicator.Position2

	// Set to max value
	knob.SetValue(100)
	renderer.Layout(size)
	maxPos := renderer.indicator.Position2

	// Positions should be different
	assert.NotEqual(t, minPos, maxPos)
}

func TestRotatingKnobRenderer_ThumbSize(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)
	size := fyne.NewSize(100, 100)

	// Normal state
	knob.hovered = false
	renderer.Layout(size)
	normalSize := renderer.thumb.Size()

	// Hovered state
	knob.hovered = true
	renderer.Layout(size)
	hoveredSize := renderer.thumb.Size()

	// Hovered thumb should be larger
	assert.Greater(t, hoveredSize.Width, normalSize.Width)
	assert.Greater(t, hoveredSize.Height, normalSize.Height)
}

func TestRotatingKnobRenderer_TickPositions(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	knob.ShowTicks = true
	knob.TickCount = 5

	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)
	size := fyne.NewSize(100, 100)
	renderer.Layout(size)

	// Check that ticks are positioned
	assert.Equal(t, 5, len(renderer.ticks))
	for _, tick := range renderer.ticks {
		// Tick positions should be set
		assert.NotEqual(t, fyne.NewPos(0, 0), tick.Position1)
		assert.NotEqual(t, fyne.NewPos(0, 0), tick.Position2)
		// Position1 and Position2 should be different
		assert.NotEqual(t, tick.Position1, tick.Position2)
	}
}

func TestRotatingKnobRenderer_Destroy(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)

	// Should not panic
	assert.NotPanics(t, func() {
		renderer.Destroy()
	})
}

func TestRotatingKnob_GetAngleFromPoint(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	knob.Resize(fyne.NewSize(100, 100))

	// Point at top (north) should be around 0°
	angle := knob.getAngleFromPoint(fyne.NewPos(50, 0))
	// Due to border cases and rounding, we check it's in the right region
	assert.True(t, angle < 10.0 || angle > 350.0)

	// Point at right (east) should be around 90°
	angle = knob.getAngleFromPoint(fyne.NewPos(100, 50))
	assert.InDelta(t, 90.0, angle, 10.0)

	// Point at bottom (south) should be around 180°
	angle = knob.getAngleFromPoint(fyne.NewPos(50, 100))
	assert.InDelta(t, 180.0, angle, 10.0)

	// Point at left (west) should be around 270°
	angle = knob.getAngleFromPoint(fyne.NewPos(0, 50))
	assert.InDelta(t, 270.0, angle, 10.0)
}

func TestRotatingKnob_UpdateValueFromAngle(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	knob.StartAngle = 0   // Top
	knob.EndAngle = 180   // Bottom

	// Angle at start should give min value
	knob.updateValueFromAngle(0)
	assert.InDelta(t, 0.0, knob.Value, 1.0)

	// Angle at end should give max value
	knob.updateValueFromAngle(180)
	assert.InDelta(t, 100.0, knob.Value, 1.0)

	// Angle at midpoint should give mid value
	knob.updateValueFromAngle(90)
	assert.InDelta(t, 50.0, knob.Value, 1.0)
}

func TestRotatingKnob_UpdateValueFromAngleWrapping(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	knob.Wrapping = true
	knob.StartAngle = 0
	knob.EndAngle = 180

	// Angle beyond end should wrap when wrapping is enabled
	knob.updateValueFromAngle(270)
	assert.GreaterOrEqual(t, knob.Value, 0.0)
	assert.LessOrEqual(t, knob.Value, 100.0)
}

func TestRotatingKnob_ColorConsistency(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)

	// Check colors are set properly
	renderer.Refresh()

	// Track should be disabled color
	assert.Equal(t, theme.DisabledColor(), renderer.track.StrokeColor)

	// Active should be primary color
	assert.Equal(t, theme.PrimaryColor(), renderer.active.StrokeColor)

	// Indicator should be foreground color
	assert.Equal(t, theme.ForegroundColor(), renderer.indicator.StrokeColor)
}

func TestRotatingKnob_RectangularSize(t *testing.T) {
	knob := NewRotatingKnob(0, 100)
	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)

	// Non-square size
	size := fyne.NewSize(120, 80)
	renderer.Layout(size)

	// Should use the smaller dimension for diameter
	expectedDiameter := float32(80) * 0.85
	assert.Equal(t, expectedDiameter, renderer.track.Size().Width)
	assert.Equal(t, expectedDiameter, renderer.track.Size().Height)
}

func TestRotatingKnob_ExtendedRange(t *testing.T) {
	knob := NewRotatingKnob(0, 1000)
	knob.StartAngle = -180 // Full circle minus a gap
	knob.EndAngle = 170

	knob.SetValue(500)
	assert.Equal(t, 500.0, knob.Value)

	renderer := test.TempWidgetRenderer(t, knob).(*rotatingKnobRenderer)
	size := fyne.NewSize(100, 100)
	renderer.Layout(size)

	// Should position indicator correctly for mid-value
	assert.NotEqual(t, renderer.indicator.Position1, renderer.indicator.Position2)
}

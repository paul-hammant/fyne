package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
)

type extendedRotaryControl struct {
	RotaryControl
}

func newExtendedRotaryControl() *extendedRotaryControl {
	knob := &extendedRotaryControl{}
	knob.ExtendBaseWidget(knob)
	knob.Min = 0
	knob.Max = 100
	knob.Value = 50
	knob.StartAngle = -135
	knob.EndAngle = 135
	return knob
}

func TestRotaryControl_Extended_Value(t *testing.T) {
	knob := newExtendedRotaryControl()
	knob.Resize(knob.MinSize())
	objs := cache.Renderer(knob).Objects()
	assert.GreaterOrEqual(t, len(objs), 4) // track, active, indicator, thumb, centerDot, optionally ticks

	// Get thumb position at value 50
	thumb := objs[3] // thumb is 4th object (after track, active, indicator)
	thumbPos := thumb.Position()

	// Change value and verify thumb moved
	knob.Value = 75
	knob.Refresh()
	assert.NotEqual(t, thumbPos, thumb.Position())
}

func TestRotaryControl_Extended_Drag(t *testing.T) {
	knob := newExtendedRotaryControl()
	knob.Resize(fyne.NewSize(100, 100))
	objs := cache.Renderer(knob).Objects()
	assert.GreaterOrEqual(t, len(objs), 4)

	initialValue := knob.Value

	// Drag to a new position
	drag := &fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(80, 50)}}
	knob.Dragged(drag)
	assert.NotEqual(t, initialValue, knob.Value)
}

func TestRotaryControl_Extended_MinSize(t *testing.T) {
	knob := newExtendedRotaryControl()
	minSize := knob.MinSize()

	// Should have reasonable minimum size
	assert.Greater(t, minSize.Width, float32(0))
	assert.Greater(t, minSize.Height, float32(0))
	// Rotary controls should be square
	assert.Equal(t, minSize.Width, minSize.Height)
}

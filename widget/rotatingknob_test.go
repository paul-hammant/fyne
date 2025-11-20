package widget_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestRotatingKnob_Creation(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	assert.NotNil(t, knob)
	assert.Equal(t, 0.0, knob.Min)
	assert.Equal(t, 100.0, knob.Max)
	assert.Equal(t, 50.0, knob.Value) // Default to midpoint
	assert.False(t, knob.Disabled())
}

func TestRotatingKnob_SetValue(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	// Set normal value
	knob.SetValue(75)
	assert.Equal(t, 75.0, knob.Value)

	// Set to min
	knob.SetValue(0)
	assert.Equal(t, 0.0, knob.Value)

	// Set to max
	knob.SetValue(100)
	assert.Equal(t, 100.0, knob.Value)
}

func TestRotatingKnob_SetValueClamping(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	// Value below min should be clamped
	knob.SetValue(-10)
	assert.Equal(t, 0.0, knob.Value)

	// Value above max should be clamped
	knob.SetValue(150)
	assert.Equal(t, 100.0, knob.Value)
}

func TestRotatingKnob_SetValueWrapping(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.Wrapping = true

	// Value below min should wrap
	knob.SetValue(-10)
	assert.Equal(t, 90.0, knob.Value)

	// Value above max should wrap
	knob.SetValue(110)
	assert.Equal(t, 10.0, knob.Value)
}

func TestRotatingKnob_OnChanged(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	var changedValue float64
	changedCalled := false
	knob.OnChanged = func(value float64) {
		changedValue = value
		changedCalled = true
	}

	knob.SetValue(75)

	assert.True(t, changedCalled)
	assert.Equal(t, 75.0, changedValue)
}

func TestRotatingKnob_OnChangeEnded(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	var endedValue float64
	endedCalled := false
	knob.OnChangeEnded = func(value float64) {
		endedValue = value
		endedCalled = true
	}

	// Tap should trigger OnChangeEnded
	test.Tap(knob)

	assert.True(t, endedCalled)
	assert.Equal(t, knob.Value, endedValue)
}

func TestRotatingKnob_Disable(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	assert.False(t, knob.Disabled())

	knob.Disable()
	assert.True(t, knob.Disabled())

	knob.Enable()
	assert.False(t, knob.Disabled())
}

func TestRotatingKnob_DisabledInteraction(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.SetValue(50)
	knob.Disable()

	changedCalled := false
	knob.OnChanged = func(value float64) {
		changedCalled = true
	}

	// Drag event should be ignored when disabled
	drag := &fyne.DragEvent{}
	drag.PointEvent.Position = fyne.NewPos(30, 30)
	knob.Dragged(drag)
	assert.False(t, changedCalled)
	assert.Equal(t, 50.0, knob.Value) // Value unchanged
}

func TestRotatingKnob_KeyboardInput(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.SetValue(50)
	knob.Step = 10

	// Up arrow should increase
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	assert.Equal(t, 60.0, knob.Value)

	// Right arrow should increase
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Equal(t, 70.0, knob.Value)

	// Down arrow should decrease
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Equal(t, 60.0, knob.Value)

	// Left arrow should decrease
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.Equal(t, 50.0, knob.Value)
}

func TestRotatingKnob_KeyboardInputHome(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.SetValue(50)

	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyHome})
	assert.Equal(t, 0.0, knob.Value)
}

func TestRotatingKnob_KeyboardInputEnd(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.SetValue(50)

	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnd})
	assert.Equal(t, 100.0, knob.Value)
}

func TestRotatingKnob_KeyboardPageUpDown(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.SetValue(50)
	knob.Step = 1

	// Page up should increase by 10x step
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyPageUp})
	assert.Equal(t, 60.0, knob.Value)

	// Page down should decrease by 10x step
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyPageDown})
	assert.Equal(t, 50.0, knob.Value)
}

func TestRotatingKnob_Tapped(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.SetValue(50)

	test.Tap(knob)

	// Value should be within valid range after tap
	assert.True(t, knob.Value >= knob.Min && knob.Value <= knob.Max)
}

func TestRotatingKnob_Dragged(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.Resize(fyne.NewSize(100, 100))
	knob.SetValue(50)

	// Store initial value
	initialValue := knob.Value

	// Drag in a circular motion
	drag := &fyne.DragEvent{}
	drag.PointEvent.Position = fyne.NewPos(80, 20)
	knob.Dragged(drag)

	// Value should have changed
	assert.NotEqual(t, initialValue, knob.Value)
	assert.True(t, knob.Value >= knob.Min && knob.Value <= knob.Max)
}

func TestRotatingKnob_DragEnd(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	endedCalled := false
	knob.OnChangeEnded = func(value float64) {
		endedCalled = true
	}

	// Simulate drag
	knob.Dragged(&fyne.DragEvent{
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(50, 10)},
	})
	knob.DragEnd()

	assert.True(t, endedCalled)
}

func TestRotatingKnob_Scrolled(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.SetValue(50)
	knob.Step = 5

	// Scroll up should increase
	knob.Scrolled(&fyne.ScrollEvent{
		Scrolled: fyne.NewDelta(0, 1),
	})
	assert.Equal(t, 55.0, knob.Value)

	// Scroll down should decrease
	knob.Scrolled(&fyne.ScrollEvent{
		Scrolled: fyne.NewDelta(0, -1),
	})
	assert.Equal(t, 50.0, knob.Value)
}

func TestRotatingKnob_FocusGainedLost(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	// Gain focus
	knob.FocusGained()
	// Focus state is internal

	// Lose focus
	knob.FocusLost()
	// Focus state is internal
}

func TestRotatingKnob_MinSize(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	minSize := knob.MinSize()

	assert.Greater(t, minSize.Width, float32(0))
	assert.Greater(t, minSize.Height, float32(0))
	assert.Equal(t, minSize.Width, minSize.Height) // Should be square
}

func TestRotatingKnob_Renderer(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	renderer := test.TempWidgetRenderer(t, knob)

	assert.NotNil(t, renderer)
	assert.Greater(t, len(renderer.Objects()), 0)
}

func TestRotatingKnob_StartEndAngles(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	// Default angles
	assert.Equal(t, -135.0, knob.StartAngle)
	assert.Equal(t, 135.0, knob.EndAngle)

	// Custom angles
	knob.StartAngle = 0
	knob.EndAngle = 180
	assert.Equal(t, 0.0, knob.StartAngle)
	assert.Equal(t, 180.0, knob.EndAngle)
}

func TestRotatingKnob_Ticks(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	// Ticks enabled by default
	assert.True(t, knob.ShowTicks)
	assert.Equal(t, 11, knob.TickCount)

	// Disable ticks
	knob.ShowTicks = false
	assert.False(t, knob.ShowTicks)
}

func TestRotatingKnob_DataBinding(t *testing.T) {
	val := binding.NewFloat()
	val.Set(75.0)

	knob := widget.NewRotatingKnobWithData(0, 100, val)

	// Initial value from binding
	assert.Equal(t, 75.0, knob.Value)

	// Change knob value, binding should update
	knob.SetValue(50)
	boundValue, err := val.Get()
	assert.NoError(t, err)
	assert.Equal(t, 50.0, boundValue)

	// Change binding value, knob should update
	val.Set(25)
	assert.Equal(t, 25.0, knob.Value)
}

func TestRotatingKnob_NegativeRange(t *testing.T) {
	knob := widget.NewRotatingKnob(-50, 50)

	assert.Equal(t, -50.0, knob.Min)
	assert.Equal(t, 50.0, knob.Max)
	assert.Equal(t, 0.0, knob.Value) // Midpoint

	knob.SetValue(-25)
	assert.Equal(t, -25.0, knob.Value)

	knob.SetValue(25)
	assert.Equal(t, 25.0, knob.Value)
}

func TestRotatingKnob_ZeroStep(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.Step = 0 // Zero step should use default (1% of range)
	knob.SetValue(50)

	// Should still work with arrow keys (using default step)
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	assert.Greater(t, knob.Value, 50.0)
}

func TestRotatingKnob_SameValueNoCallback(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)
	knob.SetValue(50)

	callCount := 0
	knob.OnChanged = func(value float64) {
		callCount++
	}

	// Setting same value should not trigger callback
	knob.SetValue(50)
	assert.Equal(t, 0, callCount)

	// Setting different value should trigger callback
	knob.SetValue(51)
	assert.Equal(t, 1, callCount)
}

func TestRotatingKnob_MouseHover(t *testing.T) {
	knob := widget.NewRotatingKnob(0, 100)

	// Simulate mouse enter
	knob.MouseIn(nil)
	// No direct way to check hovered state from outside, but ensure no panic

	// Simulate mouse move
	knob.MouseMoved(nil)

	// Simulate mouse exit
	knob.MouseOut()
}

package widget_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestRotaryControl_Creation(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)

	assert.NotNil(t, knob)
	assert.Equal(t, 0.0, knob.Min)
	assert.Equal(t, 100.0, knob.Max)
	assert.Equal(t, 50.0, knob.Value) // Default to midpoint
	assert.False(t, knob.Disabled())
}

func TestRotaryControl_SetValue(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)

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

func TestRotaryControl_SetValueClamping(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)

	// Value below min should be clamped
	knob.SetValue(-10)
	assert.Equal(t, 0.0, knob.Value)

	// Value above max should be clamped
	knob.SetValue(150)
	assert.Equal(t, 100.0, knob.Value)
}

func TestRotaryControl_SetValueWrapping(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.Wrapping = true

	// Value below min should wrap
	knob.SetValue(-10)
	assert.Equal(t, 90.0, knob.Value)

	// Value above max should wrap
	knob.SetValue(110)
	assert.Equal(t, 10.0, knob.Value)
}

func TestRotaryControl_OnChangeEnded(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)

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

func TestRotaryControl_Disable(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)

	assert.False(t, knob.Disabled())

	knob.Disable()
	assert.True(t, knob.Disabled())

	knob.Enable()
	assert.False(t, knob.Disabled())
}

func TestRotaryControl_DisabledInteraction(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
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

func TestRotaryControl_KeyboardInput(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
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

func TestRotaryControl_KeyboardInputHome(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.SetValue(50)

	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyHome})
	assert.Equal(t, 0.0, knob.Value)
}

func TestRotaryControl_KeyboardInputEnd(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.SetValue(50)

	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnd})
	assert.Equal(t, 100.0, knob.Value)
}

func TestRotaryControl_KeyboardPageUpDown(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.SetValue(50)
	knob.Step = 1

	// Page up should increase by 10x step
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyPageUp})
	assert.Equal(t, 60.0, knob.Value)

	// Page down should decrease by 10x step
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyPageDown})
	assert.Equal(t, 50.0, knob.Value)
}

func TestRotaryControl_Tapped(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.SetValue(50)

	test.Tap(knob)

	// Value should be within valid range after tap
	assert.True(t, knob.Value >= knob.Min && knob.Value <= knob.Max)
}

func TestRotaryControl_Dragged(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
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

func TestRotaryControl_DragEnd(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)

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

func TestRotaryControl_Scrolled(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
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

func TestRotaryControl_FocusGainedLost(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)

	// Gain focus
	knob.FocusGained()
	// Focus state is internal

	// Lose focus
	knob.FocusLost()
	// Focus state is internal
}

func TestRotaryControl_MinSize(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	minSize := knob.MinSize()

	assert.Greater(t, minSize.Width, float32(0))
	assert.Greater(t, minSize.Height, float32(0))
	assert.Equal(t, minSize.Width, minSize.Height) // Should be square
}

func TestRotaryControl_Renderer(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	renderer := test.TempWidgetRenderer(t, knob)

	assert.NotNil(t, renderer)
	assert.Greater(t, len(renderer.Objects()), 0)
}

func TestRotaryControl_StartEndAngles(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)

	// Default angles
	assert.Equal(t, -135.0, knob.StartAngle)
	assert.Equal(t, 135.0, knob.EndAngle)

	// Custom angles
	knob.StartAngle = 0
	knob.EndAngle = 180
	assert.Equal(t, 0.0, knob.StartAngle)
	assert.Equal(t, 180.0, knob.EndAngle)
}

func TestRotaryControl_Ticks(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)

	// Ticks enabled by default
	assert.True(t, knob.ShowTicks)
	assert.Equal(t, 11, knob.TickCount)

	// Disable ticks
	knob.ShowTicks = false
	assert.False(t, knob.ShowTicks)
}

func TestRotaryControl_DataBinding(t *testing.T) {
	val := binding.NewFloat()
	val.Set(75.0)

	knob := widget.NewRotaryControlWithData(0, 100, val)

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

func TestRotaryControl_NegativeRange(t *testing.T) {
	knob := widget.NewRotaryControl(-50, 50)

	assert.Equal(t, -50.0, knob.Min)
	assert.Equal(t, 50.0, knob.Max)
	assert.Equal(t, 0.0, knob.Value) // Midpoint

	knob.SetValue(-25)
	assert.Equal(t, -25.0, knob.Value)

	knob.SetValue(25)
	assert.Equal(t, 25.0, knob.Value)
}

func TestRotaryControl_ZeroStep(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.Step = 0 // Zero step should use default (1% of range)
	knob.SetValue(50)

	// Should still work with arrow keys (using default step)
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	assert.Greater(t, knob.Value, 50.0)
}

func TestRotaryControl_SameValueNoCallback(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
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

func TestRotaryControl_MouseHover(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)

	// Simulate mouse enter
	knob.MouseIn(nil)
	// No direct way to check hovered state from outside, but ensure no panic

	// Simulate mouse move
	knob.MouseMoved(nil)

	// Simulate mouse exit
	knob.MouseOut()
}

func TestRotaryControl_Binding(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.SetValue(20)
	assert.Equal(t, 20.0, knob.Value)

	val := binding.NewFloat()
	knob.Bind(val)
	test.WidgetRenderer(knob) // force render to process binding
	assert.Equal(t, 0.0, knob.Value)

	err := val.Set(30)
	assert.NoError(t, err)
	test.WidgetRenderer(knob) // force render to process binding
	assert.Equal(t, 30.0, knob.Value)

	knob.SetValue(50)
	f, err := val.Get()
	assert.NoError(t, err)
	assert.Equal(t, 50.0, f)

	knob.Unbind()
	assert.Equal(t, 50.0, knob.Value)
}

func TestRotaryControl_OnChangedComprehensive(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.Resize(fyne.NewSize(100, 100))
	assert.Nil(t, knob.OnChanged)

	changes := 0
	knob.OnChanged = func(_ float64) {
		changes++
	}

	assert.Equal(t, 0, changes)

	// SetValue should trigger OnChanged
	knob.SetValue(25)
	assert.Equal(t, 1, changes)

	// Drag should trigger OnChanged
	drag := &fyne.DragEvent{}
	drag.PointEvent.Position = fyne.NewPos(80, 50)
	knob.Dragged(drag)
	assert.Equal(t, 2, changes)

	// Same position drag should not trigger (no value change)
	knob.Dragged(drag)
	assert.Equal(t, 2, changes)

	// Different position should trigger
	drag.PointEvent.Position = fyne.NewPos(50, 20)
	knob.Dragged(drag)
	assert.Equal(t, 3, changes)

	// Tap should trigger
	tap := &fyne.PointEvent{}
	tap.Position = fyne.NewPos(20, 50)
	knob.Tapped(tap)
	assert.Equal(t, 4, changes)

	// Key should trigger
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Equal(t, 5, changes)

	// Scroll should trigger
	knob.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, 1)})
	assert.Equal(t, 6, changes)
}

func TestRotaryControl_OnChangeEndedComprehensive(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.Resize(fyne.NewSize(100, 100))
	assert.Nil(t, knob.OnChangeEnded)

	changes := 0
	knob.OnChangeEnded = func(_ float64) {
		changes++
	}

	assert.Equal(t, 0, changes)

	// SetValue should NOT trigger OnChangeEnded
	knob.SetValue(25)
	assert.Equal(t, 0, changes)

	// Drag should NOT trigger OnChangeEnded (only DragEnd does)
	drag := &fyne.DragEvent{}
	drag.PointEvent.Position = fyne.NewPos(80, 50)
	knob.Dragged(drag)
	assert.Equal(t, 0, changes)

	// DragEnd should trigger OnChangeEnded
	knob.DragEnd()
	assert.Equal(t, 1, changes)

	// Tap should trigger OnChangeEnded
	tap := &fyne.PointEvent{}
	tap.Position = fyne.NewPos(20, 50)
	knob.Tapped(tap)
	assert.Equal(t, 2, changes)

	// Same position tap should NOT trigger (no value change)
	knob.Tapped(tap)
	assert.Equal(t, 2, changes)

	// Different position tap should trigger
	tap.Position = fyne.NewPos(50, 80)
	knob.Tapped(tap)
	assert.Equal(t, 3, changes)

	// Key should trigger OnChangeEnded
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Equal(t, 4, changes)
}

func TestRotaryControl_FocusDesktop(t *testing.T) {
	if fyne.CurrentDevice().IsMobile() {
		return
	}
	knob := widget.NewRotaryControl(0, 100)
	win := test.NewTempWindow(t, knob)
	test.Tap(knob)

	assert.Equal(t, win.Canvas().Focused(), knob)
}

func TestRotaryControl_FocusComprehensive(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.Step = 10
	knob.SetValue(0)

	knob.FocusGained()
	// Focus state is internal, just verify no panic

	knob.FocusLost()
	// Focus state is internal, just verify no panic

	knob.MouseIn(nil)
	// Hovered state is internal, just verify no panic

	knob.MouseOut()
	// Hovered state is internal, just verify no panic

	// Test keyboard iteration - Up/Right increase
	for i := 1; i <= 10; i++ {
		knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, float64(i*10), knob.Value)
	}

	// Should stay at max
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Equal(t, knob.Max, knob.Value)

	// Test keyboard iteration - Down/Left decrease
	for i := 9; i >= 0; i-- {
		knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, float64(i*10), knob.Value)
	}

	// Should stay at min
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.Equal(t, knob.Min, knob.Value)
}

func TestRotaryControl_DisabledComprehensive(t *testing.T) {
	knob := widget.NewRotaryControl(0, 100)
	knob.Resize(fyne.NewSize(100, 100))
	knob.SetValue(0) // Start at min
	knob.Disable()

	changes := 0
	knob.OnChanged = func(_ float64) {
		changes++
	}

	// Tap should be ignored when disabled
	tap := &fyne.PointEvent{}
	tap.Position = fyne.NewPos(80, 50)
	knob.Tapped(tap)
	assert.Equal(t, 0, changes)

	// Drag should be ignored when disabled
	drag := &fyne.DragEvent{}
	drag.PointEvent.Position = fyne.NewPos(50, 20) // Top center = ~middle value
	knob.Dragged(drag)
	assert.Equal(t, 0, changes)

	// Key should be ignored when disabled
	knob.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Equal(t, 0, changes)

	// Scroll should be ignored when disabled
	knob.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, 1)})
	assert.Equal(t, 0, changes)

	// Enable and verify interaction works
	knob.Enable()
	knob.Dragged(drag) // Should change from 0 to ~middle value
	assert.Equal(t, 1, changes)
}

func TestRotaryControl_Layout(t *testing.T) {
	test.NewTempApp(t)

	knob := widget.NewRotaryControl(0, 100)
	knob.Resize(fyne.NewSize(100, 100))

	w := test.NewWindow(knob)
	defer w.Close()
	w.Resize(fyne.NewSize(120, 120))

	test.AssertRendersToMarkup(t, "rotarycontrol/layout.xml", w.Canvas())
}

func TestRotaryControl_LayoutDisabled(t *testing.T) {
	test.NewTempApp(t)

	knob := widget.NewRotaryControl(0, 100)
	knob.Resize(fyne.NewSize(100, 100))
	knob.Disable()

	w := test.NewWindow(knob)
	defer w.Close()
	w.Resize(fyne.NewSize(120, 120))

	test.AssertRendersToMarkup(t, "rotarycontrol/layout_disabled.xml", w.Canvas())
}

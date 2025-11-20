package widget

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

// Declare compile-time interface conformance
var _ fyne.Widget = (*RotatingKnob)(nil)
var _ fyne.Draggable = (*RotatingKnob)(nil)
var _ fyne.Tappable = (*RotatingKnob)(nil)
var _ fyne.Focusable = (*RotatingKnob)(nil)
var _ desktop.Hoverable = (*RotatingKnob)(nil)
var _ fyne.Disableable = (*RotatingKnob)(nil)

// RotatingKnob is a widget that provides a circular dial/knob control for selecting values
// within a range, similar to a potentiometer or volume knob.
//
// The knob can be controlled via:
// - Mouse/touch dragging (circular rotation)
// - Clicking/tapping at a position
// - Keyboard arrow keys (when focused)
// - Scroll wheel (when hovered)
//
// Example usage:
//
//	knob := widget.NewRotatingKnob(0, 100)
//	knob.SetValue(50)
//	knob.OnChanged = func(value float64) {
//	    fmt.Printf("Value changed to: %.2f\n", value)
//	}
type RotatingKnob struct {
	DisableableWidget

	// Value is the current value of the knob
	Value float64
	// Min is the minimum value
	Min float64
	// Max is the maximum value
	Max float64
	// Step is the increment for keyboard adjustments (0 for continuous)
	Step float64

	// StartAngle is the angle in degrees where the knob range starts (0° = top, clockwise)
	// Default is -135° (bottom-left)
	StartAngle float64
	// EndAngle is the angle in degrees where the knob range ends (0° = top, clockwise)
	// Default is 135° (bottom-right)
	EndAngle float64

	// Wrapping enables wrapping from max back to min (and vice versa)
	Wrapping bool
	// ShowTicks enables visual tick marks around the knob
	ShowTicks bool
	// TickCount is the number of tick marks to show (if ShowTicks is true)
	TickCount int

	// AccentColor is the color used for the active arc and thumb (nil uses theme color)
	AccentColor color.Color
	// TrackColor is the color used for the background track (nil uses theme color)
	TrackColor color.Color
	// WedgeColor is the color used for the wedge backdrop fill (nil disables wedge)
	WedgeColor color.Color

	// OnChanged is called when the value changes (during dragging)
	OnChanged func(float64)
	// OnChangeEnded is called when a value change ends (drag end, key release)
	OnChangeEnded func(float64)

	binder  basicBinder
	hovered bool
	focused bool
}

// NewRotatingKnob creates a new rotating knob widget with the specified min and max values.
// The knob is initialized with a value at the midpoint of the range.
func NewRotatingKnob(min, max float64) *RotatingKnob {
	knob := &RotatingKnob{
		Value:      (min + max) / 2,
		Min:        min,
		Max:        max,
		Step:       (max - min) / 100, // Default to 1% of range
		StartAngle: -135,              // Bottom-left
		EndAngle:   135,                // Bottom-right (270° total sweep)
		ShowTicks:  true,
		TickCount:  11, // 0, 10, 20, ... 100 for percentage-like display
	}
	knob.ExtendBaseWidget(knob)
	return knob
}

// NewRotatingKnobWithData creates a new rotating knob bound to a float data item.
//
// Since: 2.0
func NewRotatingKnobWithData(min, max float64, data binding.Float) *RotatingKnob {
	knob := NewRotatingKnob(min, max)
	knob.Bind(data)
	return knob
}

// Bind connects the specified data source to this RotatingKnob.
// The current value will be displayed and any changes in the data will cause the widget to update.
// User interactions with this RotatingKnob will set the value into the data source.
//
// Since: 2.0
func (k *RotatingKnob) Bind(data binding.Float) {
	k.binder.SetCallback(k.updateFromData)
	k.binder.Bind(data)

	k.OnChanged = func(_ float64) {
		k.binder.CallWithData(k.writeData)
	}
}

// Unbind disconnects any configured data source from this RotatingKnob.
// The current value will remain at the last value of the data source.
//
// Since: 2.0
func (k *RotatingKnob) Unbind() {
	k.OnChanged = nil
	k.binder.Unbind()
}

// updateFromData is called when the data changes
func (k *RotatingKnob) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	floatSource, ok := data.(binding.Float)
	if !ok {
		return
	}
	val, err := floatSource.Get()
	if err != nil {
		return
	}
	k.SetValue(val)
}

// writeData writes the current value to the data binding
func (k *RotatingKnob) writeData(data binding.DataItem) {
	if data == nil {
		return
	}
	floatTarget, ok := data.(binding.Float)
	if !ok {
		return
	}
	floatTarget.Set(k.Value)
}

// SetValue updates the knob value and refreshes the widget
func (k *RotatingKnob) SetValue(value float64) {
	// Clamp to range (unless wrapping)
	if !k.Wrapping {
		if value < k.Min {
			value = k.Min
		}
		if value > k.Max {
			value = k.Max
		}
	} else {
		// Wrap around
		valueRange := k.Max - k.Min
		for value < k.Min {
			value += valueRange
		}
		for value > k.Max {
			value -= valueRange
		}
	}

	if k.Value == value {
		return
	}

	k.Value = value
	k.Refresh()

	if k.OnChanged != nil {
		k.OnChanged(k.Value)
	}
}

// MinSize returns the minimum size for the knob
func (k *RotatingKnob) MinSize() fyne.Size {
	k.ExtendBaseWidget(k)
	return k.BaseWidget.MinSize()
}

// CreateRenderer creates the renderer for the rotating knob
func (k *RotatingKnob) CreateRenderer() fyne.WidgetRenderer {
	k.ExtendBaseWidget(k)

	// Wedge backdrop (thick arc showing current value range)
	var wedge *canvas.Arc
	if k.WedgeColor != nil {
		wedge = canvas.NewArc(0, 0, 0.9, color.Transparent) // High cutout ratio to minimize inner edge
		wedge.StrokeColor = k.WedgeColor
		wedge.StrokeWidth = 20 // Thick stroke
	}

	// Track arc (the full range available)
	track := canvas.NewArc(0, 0, 0.9, color.Transparent) // High cutout ratio to minimize inner edge visibility
	track.StrokeWidth = 8
	track.StrokeColor = theme.DisabledColor()

	// Active arc (the current value indicator)
	active := canvas.NewArc(0, 0, 0.9, color.Transparent) // High cutout ratio to minimize inner edge visibility
	active.StrokeWidth = 8
	active.StrokeColor = theme.PrimaryColor()

	// Indicator line (points to current value)
	indicator := canvas.NewLine(theme.ForegroundColor())
	indicator.StrokeWidth = 3

	// Thumb (circle at the indicator tip)
	thumb := canvas.NewCircle(theme.ForegroundColor())

	// Center dot
	centerDot := canvas.NewCircle(theme.BackgroundColor())

	objects := []fyne.CanvasObject{}
	if wedge != nil {
		objects = append(objects, wedge)
	}
	objects = append(objects, track, active, indicator, thumb, centerDot)

	// Add tick marks if enabled
	var ticks []*canvas.Line
	if k.ShowTicks && k.TickCount > 0 {
		for i := 0; i < k.TickCount; i++ {
			tick := canvas.NewLine(theme.DisabledColor())
			tick.StrokeWidth = 1
			ticks = append(ticks, tick)
			objects = append(objects, tick)
		}
	}

	r := &rotatingKnobRenderer{
		knob:      k,
		wedge:     wedge,
		track:     track,
		active:    active,
		indicator: indicator,
		thumb:     thumb,
		centerDot: centerDot,
		ticks:     ticks,
	}
	r.objects = objects
	r.Refresh()
	return r
}

// Dragged handles drag events for rotating the knob
func (k *RotatingKnob) Dragged(e *fyne.DragEvent) {
	if k.Disabled() {
		return
	}

	angle := k.getAngleFromPoint(e.Position)
	k.updateValueFromAngle(angle)
}

// DragEnd is called when dragging ends
func (k *RotatingKnob) DragEnd() {
	if k.OnChangeEnded != nil {
		k.OnChangeEnded(k.Value)
	}
}

// Tapped handles tap events for jumping to a position
func (k *RotatingKnob) Tapped(e *fyne.PointEvent) {
	if k.Disabled() {
		return
	}

	angle := k.getAngleFromPoint(e.Position)
	k.updateValueFromAngle(angle)

	if k.OnChangeEnded != nil {
		k.OnChangeEnded(k.Value)
	}
}

// FocusGained is called when the knob gains focus
func (k *RotatingKnob) FocusGained() {
	k.focused = true
	k.Refresh()
}

// FocusLost is called when the knob loses focus
func (k *RotatingKnob) FocusLost() {
	k.focused = false
	k.Refresh()
}

// TypedRune handles rune input (not used for knob)
func (k *RotatingKnob) TypedRune(_ rune) {
	// Not used
}

// TypedKey handles keyboard input for adjusting the knob value
func (k *RotatingKnob) TypedKey(key *fyne.KeyEvent) {
	if k.Disabled() {
		return
	}

	step := k.Step
	if step == 0 {
		step = (k.Max - k.Min) / 100
	}

	switch key.Name {
	case fyne.KeyUp, fyne.KeyRight:
		k.SetValue(k.Value + step)
		if k.OnChangeEnded != nil {
			k.OnChangeEnded(k.Value)
		}
	case fyne.KeyDown, fyne.KeyLeft:
		k.SetValue(k.Value - step)
		if k.OnChangeEnded != nil {
			k.OnChangeEnded(k.Value)
		}
	case fyne.KeyPageUp:
		k.SetValue(k.Value + step*10)
		if k.OnChangeEnded != nil {
			k.OnChangeEnded(k.Value)
		}
	case fyne.KeyPageDown:
		k.SetValue(k.Value - step*10)
		if k.OnChangeEnded != nil {
			k.OnChangeEnded(k.Value)
		}
	case fyne.KeyHome:
		k.SetValue(k.Min)
		if k.OnChangeEnded != nil {
			k.OnChangeEnded(k.Value)
		}
	case fyne.KeyEnd:
		k.SetValue(k.Max)
		if k.OnChangeEnded != nil {
			k.OnChangeEnded(k.Value)
		}
	}
}

// MouseIn handles mouse enter events
func (k *RotatingKnob) MouseIn(_ *desktop.MouseEvent) {
	k.hovered = true
	k.Refresh()
}

// MouseMoved handles mouse move events
func (k *RotatingKnob) MouseMoved(_ *desktop.MouseEvent) {
	// Visual feedback could be added here
}

// MouseOut handles mouse exit events
func (k *RotatingKnob) MouseOut() {
	k.hovered = false
	k.Refresh()
}

// Scrolled handles scroll wheel events for adjusting the value
func (k *RotatingKnob) Scrolled(e *fyne.ScrollEvent) {
	if k.Disabled() {
		return
	}

	step := k.Step
	if step == 0 {
		step = (k.Max - k.Min) / 100
	}

	// Scroll up increases value, scroll down decreases
	if e.Scrolled.DY > 0 {
		k.SetValue(k.Value + step)
	} else if e.Scrolled.DY < 0 {
		k.SetValue(k.Value - step)
	}

	if k.OnChangeEnded != nil {
		k.OnChangeEnded(k.Value)
	}
}

// getAngleFromPoint calculates the angle in degrees from a point relative to the knob center
func (k *RotatingKnob) getAngleFromPoint(pos fyne.Position) float64 {
	size := k.Size()
	centerX := size.Width / 2
	centerY := size.Height / 2

	dx := pos.X - centerX
	dy := pos.Y - centerY

	// Calculate angle using atan2 (returns radians, -π to π)
	// For 0° at top (north) and clockwise positive, we use atan2(dx, -dy)
	radians := math.Atan2(float64(dx), float64(-dy))

	// Convert to degrees (0-360)
	degrees := radians * 180 / math.Pi

	// Normalize to 0-360
	if degrees < 0 {
		degrees += 360
	}

	return degrees
}

// updateValueFromAngle updates the knob value based on an angle
func (k *RotatingKnob) updateValueFromAngle(angle float64) {
	// Normalize start and end angles
	startAngle := k.StartAngle
	endAngle := k.EndAngle

	// Normalize to 0-360
	for startAngle < 0 {
		startAngle += 360
	}
	for endAngle < 0 {
		endAngle += 360
	}
	for startAngle >= 360 {
		startAngle -= 360
	}
	for endAngle >= 360 {
		endAngle -= 360
	}

	// Calculate the sweep (angular range)
	sweep := endAngle - startAngle
	if sweep < 0 {
		sweep += 360
	}

	// Calculate angle relative to start
	relativeAngle := angle - startAngle
	if relativeAngle < 0 {
		relativeAngle += 360
	}

	// If we're wrapping, the angle is always valid
	// Otherwise, clamp to the sweep range
	if !k.Wrapping && relativeAngle > sweep {
		// We're in the dead zone - determine which boundary is closer
		deadZone := 360 - sweep

		// If in first half of dead zone (closer to max), stay at max
		// If in second half of dead zone (closer to min), stay at min
		if relativeAngle < sweep+deadZone/2 {
			relativeAngle = sweep // Stay at max
		} else {
			relativeAngle = 0 // Stay at min
		}
	}

	// Convert angle to value ratio (0.0 to 1.0)
	ratio := relativeAngle / sweep
	if ratio > 1.0 {
		ratio = math.Mod(ratio, 1.0)
	}

	// Calculate value from ratio
	value := k.Min + ratio*(k.Max-k.Min)
	k.SetValue(value)
}

// rotatingKnobRenderer is the renderer for RotatingKnob
type rotatingKnobRenderer struct {
	knob      *RotatingKnob
	wedge     *canvas.Arc
	track     *canvas.Arc
	active    *canvas.Arc
	indicator *canvas.Line
	thumb     *canvas.Circle
	centerDot *canvas.Circle
	ticks     []*canvas.Line
	objects   []fyne.CanvasObject
}

func (r *rotatingKnobRenderer) Layout(size fyne.Size) {
	diameter := fyne.Min(size.Width, size.Height)
	centerX := size.Width / 2
	centerY := size.Height / 2
	radius := diameter / 2

	// Calculate current angle for wedge
	ratio := (r.knob.Value - r.knob.Min) / (r.knob.Max - r.knob.Min)
	if r.knob.Max == r.knob.Min {
		ratio = 0
	}
	startAngle := r.knob.StartAngle
	endAngle := r.knob.EndAngle
	sweep := endAngle - startAngle
	if sweep <= 0 { // allow for wrapping
		sweep += 360
	}
	currentAngle := startAngle + ratio*sweep

	// Wedge backdrop - thick arc along circumference (same size as track/active)
	if r.wedge != nil {
		wedgeDiameter := diameter * 0.85 // Same size as track/active arcs
		wedgeRadius := wedgeDiameter / 2
		r.wedge.Resize(fyne.NewSize(wedgeDiameter, wedgeDiameter))
		r.wedge.Move(fyne.NewPos(centerX-wedgeRadius, centerY-wedgeRadius))

		// Normalize angles to 0-360 range for consistent Arc rendering
		normalizedStart := startAngle
		for normalizedStart < 0 {
			normalizedStart += 360
		}
		normalizedCurrent := currentAngle
		for normalizedCurrent < 0 {
			normalizedCurrent += 360
		}

		r.wedge.StartAngle = float32(normalizedStart)
		r.wedge.EndAngle = float32(normalizedCurrent)
	}

	// Arcs - slightly smaller ring
	arcDiameter := diameter * 0.85
	arcRadius := arcDiameter / 2
	r.track.Resize(fyne.NewSize(arcDiameter, arcDiameter))
	r.track.Move(fyne.NewPos(centerX-arcRadius, centerY-arcRadius))
	r.track.StartAngle = float32(startAngle)
	r.track.EndAngle = float32(endAngle)

	r.active.Resize(fyne.NewSize(arcDiameter, arcDiameter))
	r.active.Move(fyne.NewPos(centerX-arcRadius, centerY-arcRadius))
	r.active.StartAngle = float32(startAngle)
	r.active.EndAngle = float32(currentAngle)

	// Convert to radians for calculation (0° = top = -90° in standard coords)
	angleRad := (currentAngle - 90) * math.Pi / 180

	// Indicator line from center to edge
	indicatorLength := radius * 0.5
	indicatorEndX := centerX + float32(math.Cos(float64(angleRad))*float64(indicatorLength))
	indicatorEndY := centerY + float32(math.Sin(float64(angleRad))*float64(indicatorLength))

	r.indicator.Position1 = fyne.NewPos(centerX, centerY)
	r.indicator.Position2 = fyne.NewPos(indicatorEndX, indicatorEndY)

	// Thumb at indicator tip
	thumbPosRadius := radius * 0.65
	thumbX := centerX + float32(math.Cos(float64(angleRad))*float64(thumbPosRadius))
	thumbY := centerY + float32(math.Sin(float64(angleRad))*float64(thumbPosRadius))

	thumbRadius := float32(6)
	if r.knob.hovered {
		thumbRadius = 8
	}
	r.thumb.Resize(fyne.NewSize(thumbRadius*2, thumbRadius*2))
	r.thumb.Move(fyne.NewPos(thumbX-thumbRadius, thumbY-thumbRadius))

	// Center dot
	centerDotRadius := float32(8)
	r.centerDot.Resize(fyne.NewSize(centerDotRadius*2, centerDotRadius*2))
	r.centerDot.Move(fyne.NewPos(centerX-centerDotRadius, centerY-centerDotRadius))

	// Layout tick marks
	if r.knob.ShowTicks && len(r.ticks) > 0 {
		tickOuterRadius := radius * 0.95
		tickInnerRadius := radius * 0.8

		for i, tick := range r.ticks {
			tickRatio := float64(i) / float64(len(r.ticks)-1)
			tickAngle := startAngle + tickRatio*sweep
			tickAngleRad := (tickAngle - 90) * math.Pi / 180

			x1 := centerX + float32(math.Cos(tickAngleRad)*float64(tickInnerRadius))
			y1 := centerY + float32(math.Sin(float64(tickAngleRad))*float64(tickInnerRadius))
			x2 := centerX + float32(math.Cos(tickAngleRad)*float64(tickOuterRadius))
			y2 := centerY + float32(math.Sin(float64(tickAngleRad))*float64(tickOuterRadius))

			tick.Position1 = fyne.NewPos(x1, y1)
			tick.Position2 = fyne.NewPos(x2, y2)
		}
	}
}

func (r *rotatingKnobRenderer) MinSize() fyne.Size {
	// Minimum reasonable size for a knob
	return fyne.NewSize(60, 60)
}

func (r *rotatingKnobRenderer) Refresh() {
	// Update colors based on state
	if r.knob.Disabled() {
		if r.wedge != nil {
			r.wedge.StrokeColor = theme.DisabledColor()
			r.wedge.FillColor = color.Transparent
			r.wedge.CutoutRatio = 1.0 // No closure lines
		}
		r.track.StrokeColor = theme.DisabledColor()
		r.active.StrokeColor = theme.DisabledColor()
		r.indicator.StrokeColor = theme.DisabledColor()
		r.thumb.FillColor = theme.DisabledColor()
		r.centerDot.FillColor = theme.BackgroundColor()
		for _, tick := range r.ticks {
			tick.StrokeColor = theme.DisabledColor()
		}
	} else {
		// Wedge backdrop (thick stroke, not fill)
		if r.wedge != nil && r.knob.WedgeColor != nil {
			r.wedge.StrokeColor = r.knob.WedgeColor
			r.wedge.FillColor = color.Transparent
			r.wedge.CutoutRatio = 1.0 // No closure lines
		}

		// Track shows the full range (subtle)
		trackColor := theme.DisabledColor()
		if r.knob.TrackColor != nil {
			trackColor = r.knob.TrackColor
		}
		r.track.StrokeColor = trackColor
		r.track.FillColor = color.Transparent
		if r.knob.hovered {
			r.track.StrokeWidth = 10
		} else {
			r.track.StrokeWidth = 8
		}

		// Active shows current position (prominent)
		activeColor := theme.PrimaryColor()
		if r.knob.AccentColor != nil {
			activeColor = r.knob.AccentColor
		}
		if r.knob.hovered {
			// Brighten on hover (blend with hover color if using custom color)
			if r.knob.AccentColor == nil {
				activeColor = theme.HoverColor()
			}
		}
		r.active.StrokeColor = activeColor
		r.active.FillColor = color.Transparent
		if r.knob.hovered {
			r.active.StrokeWidth = 10
		} else {
			r.active.StrokeWidth = 8
		}

		// Indicator line
	r.indicator.StrokeColor = theme.ForegroundColor()
		if r.knob.focused {
			r.indicator.StrokeWidth = 4
		} else {
			r.indicator.StrokeWidth = 3
		}

		// Thumb
	r.thumb.FillColor = activeColor

		// Center dot
	r.centerDot.FillColor = theme.BackgroundColor()
		r.centerDot.StrokeColor = theme.ShadowColor()
		r.centerDot.StrokeWidth = 1

		// Ticks
		for _, tick := range r.ticks {
			tick.StrokeColor = theme.ShadowColor()
		}
	}

	r.Layout(r.knob.Size())
	canvas.Refresh(r.knob.super())
}

func (r *rotatingKnobRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *rotatingKnobRenderer) Destroy() {}
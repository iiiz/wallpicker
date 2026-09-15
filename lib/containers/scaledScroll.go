package containers

import (
	"math"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// Tuning constants for the smoothing animation.
const (
	easeFriction = 0.20
	easeSnap     = 0.5
	easeTickHz   = 60
)

// ScaledScroll is a scroll container that applies a scale factor to incoming
// scroll deltas and eases the view to the target offset for smooth scrolling.
type ScaledScroll struct {
	*container.Scroll

	scale float32

	target    fyne.Position
	smoothing bool
	mu        sync.Mutex
}

// NewScaledScroll returns a scroll container that multiplies the incoming
// scroll delta by the given scale and eases the view toward the target
// offset for smooth scrolling.
func NewScaledScroll(direction fyne.ScrollDirection, scale float32, content fyne.CanvasObject) *ScaledScroll {
	s := &ScaledScroll{
		Scroll: &container.Scroll{
			Direction: direction,
			Content:   content,
		},
		scale: scale,
	}
	s.Scroll.OnScrolled = func(offset fyne.Position) {
		s.mu.Lock()
		s.target = offset
		s.mu.Unlock()
	}
	s.ExtendBaseWidget(s)

	go s.smoothLoop()

	return s
}

// Scrolled accumulates the scaled delta from the input device into the
// target offset and flags the smoothing loop so the view eases toward it.
// Deltas are negated to match fyne's scroll direction convention.
func (s *ScaledScroll) Scrolled(event *fyne.ScrollEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.target.X -= event.Scrolled.DX * s.scale
	s.target.Y -= event.Scrolled.DY * s.scale
	s.smoothing = true
}

func (s *ScaledScroll) smoothLoop() {
	ticker := time.NewTicker(time.Second / easeTickHz)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		if !s.smoothing {
			s.mu.Unlock()
			continue
		}
		s.mu.Unlock()

		fyne.Do(func() {
			s.animate()
		})
	}
}

// animate runs on the UI thread, clamping the target to the scrollable
// bounds and easing the current offset one step toward it. The smoothing
// flag is cleared once the view has reached the target.
func (s *ScaledScroll) animate() {
	s.mu.Lock()
	if !s.smoothing {
		s.mu.Unlock()
		return
	}

	s.clampTarget()

	offset := s.Scroll.Offset
	next := fyne.NewPos(
		s.easeAxis(offset.X, s.target.X),
		s.easeAxis(offset.Y, s.target.Y),
	)
	s.mu.Unlock()

	s.Scroll.ScrollToOffset(next)

	s.mu.Lock()
	s.smoothing = s.target != s.Scroll.Offset
	s.mu.Unlock()
}

// easeAxis moves the offset a fixed fraction of the remaining distance to the
// target each frame, snapping to the target once the gap falls below easeSnap.
func (s *ScaledScroll) easeAxis(offset, target float32) float32 {
	next := offset + (target-offset)*easeFriction
	if math.Abs(float64(target-next)) < easeSnap {
		next = target
	}
	return next
}

// clampTarget constrains the target offset to the scrollable range so the
// view cannot ease past the content bounds.
func (s *ScaledScroll) clampTarget() {
	size := s.Size()
	min := s.Content.MinSize()

	s.target.X = clamp(s.target.X, 0, fyne.Max(0, min.Width-size.Width))
	s.target.Y = clamp(s.target.Y, 0, fyne.Max(0, min.Height-size.Height))
}

// clamp bounds v to the inclusive range [lo, hi].
func clamp(v, lo, hi float32) float32 {
	return fyne.Min(fyne.Max(v, lo), hi)
}

package containers

import (
	"math"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

const (
	easeFriction = 0.20
	easeSnap     = 0.5
	easeTickHz   = 60
)

type ScaledScroll struct {
	*container.Scroll

	scale float32

	target    fyne.Position
	smoothing bool
	mu        sync.Mutex
}

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

func (s *ScaledScroll) easeAxis(offset, target float32) float32 {
	next := offset + (target-offset)*easeFriction
	if math.Abs(float64(target-next)) < easeSnap {
		next = target
	}
	return next
}

func (s *ScaledScroll) clampTarget() {
	size := s.Size()
	min := s.Content.MinSize()

	s.target.X = clamp(s.target.X, 0, fyne.Max(0, min.Width-size.Width))
	s.target.Y = clamp(s.target.Y, 0, fyne.Max(0, min.Height-size.Height))
}

func clamp(v, lo, hi float32) float32 {
	return fyne.Min(fyne.Max(v, lo), hi)
}

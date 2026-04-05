package containers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type ScaledScroll struct {
	*container.Scroll

	scale float32
}

func NewScaledScroll(direction fyne.ScrollDirection, scale float32, content fyne.CanvasObject) *ScaledScroll {
	s := &ScaledScroll{
		Scroll: &container.Scroll{
			Direction: direction,
			Content:   content,
		},
		scale: scale,
	}
	s.ExtendBaseWidget(s)

	return s
}

func (s *ScaledScroll) Scrolled(event *fyne.ScrollEvent) {
	s.Scroll.Scrolled(
		&fyne.ScrollEvent{
			PointEvent: event.PointEvent,
			Scrolled: fyne.Delta{
				DX: float32(event.Scrolled.DX) * s.scale,
				DY: float32(event.Scrolled.DY) * s.scale,
			},
		},
	)
}

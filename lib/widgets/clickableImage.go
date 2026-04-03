package widgets

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type ClickableImage struct {
	widget.BaseWidget

	image   *canvas.Image
	onClick func()
}

func NewClickableImage(image image.Image, onClickFunc func()) *ClickableImage {
	canvasImage := canvas.NewImageFromImage(image)

	ci := &ClickableImage{image: canvasImage, onClick: onClickFunc}
	ci.image.ScaleMode = canvas.ImageScaleFastest
	ci.image.FillMode = canvas.ImageFillContain
	ci.ExtendBaseWidget(ci)

	return ci
}

func (ci *ClickableImage) SetMinSize(size fyne.Size) {
	ci.image.SetMinSize(size)
}

func (ci *ClickableImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewPadded(ci.image))
}

func (ci *ClickableImage) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (ci *ClickableImage) Tapped(*fyne.PointEvent) {
	ci.onClick()
}

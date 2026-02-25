package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/alexflint/go-arg"
	"golang.org/x/image/draw"
)

type File struct {
	Path      string
	Name      string
	Extension string
}

var args struct {
	Persist bool   `arg:"-p, --persist" help:"Persist, remain after choosing a wallpaper."`
	WithExt bool   `arg:"-e, --inc-extension" help:"Include file extension as second argument to the target command or script. ie ($2)"`
	Command string `arg:"-c, --command" default:"feh --bg-fill" help:"Settings command, expects a command to run [command] [image path] [~extension~]."`
	Dir     string `arg:"positional"`
}

func main() {
	myApp := app.NewWithID("iiiz.wallpicker")
	myWindow := myApp.NewWindow("WallPicker")
	arg.MustParse(&args)

	if args.Dir == "" {
		dir, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}

		args.Dir = dir
	}

	files := getAllowedFiles(args.Dir)

	if len(files) == 0 {
		nowall := widget.NewLabelWithStyle("No wallpapers in directory.", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

		myWindow.SetContent(container.NewVBox(nowall))
		myWindow.Resize(fyne.NewSize(400, 80))

		go exitAfterDelay(5)
	} else {

		startupLabel := widget.NewLabelWithStyle("Loading Wallpapers", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		pb := widget.NewProgressBar()

		myWindow.SetContent(container.NewVBox(startupLabel, pb))
		myWindow.Resize(fyne.NewSize(600, 340))

		go loadMainContent(files, myWindow, pb)
	}

	myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}

func loadMainContent(files []File, window fyne.Window, progress *widget.ProgressBar) {
	total := len(files)
	var wg sync.WaitGroup
	done := make(chan bool, total)
	wg.Add(total)

	wpContent := container.New(layout.NewVBoxLayout())
	sc := container.NewScroll(wpContent)

	for _, file := range files {
		go loadImage(file, wpContent, &wg, done)
	}

	go func() {
		wg.Wait()
		close(done)

		// Give the ui enough time to update progress
		time.Sleep(time.Millisecond * 10)
		window.SetContent(sc)
	}()

	complete := 0
	for range done {
		complete++

		progress.SetValue(float64(complete) / float64(total))
	}
}

func loadImage(f File, c *fyne.Container, wg *sync.WaitGroup, done chan<- bool) {
	defer wg.Done()

	image := getRescaledImage(f)

	ci := NewClickableImage(image)
	ci.image.SetMinSize(fyne.NewSize(550, 200))
	ci.OnClick = func() {
		setWallpaper(f)
	}

	ci.Refresh()
	c.Add(ci)

	done <- true
}

func getRescaledImage(f File) image.Image {
	file, _ := os.Open(f.Path)
	defer file.Close()

	var original image.Image

	switch f.Extension {
	case ".jpeg", ".jpg":
		original, _ = jpeg.Decode(file)
	case ".png":
		original, _ = png.Decode(file)
	}

	scaled := image.NewRGBA(image.Rect(0, 0, original.Bounds().Max.X/8, original.Bounds().Max.Y/8))

	draw.NearestNeighbor.Scale(scaled, scaled.Rect, original, original.Bounds(), draw.Over, nil)

	return scaled
}

func setWallpaper(file File) {
	commandFields := strings.Fields(args.Command)
	cmdName := commandFields[0]
	cmdArgs := commandFields[1:]
	cmdArgs = append(cmdArgs, file.Path)

	if args.WithExt {
		cmdArgs = append(cmdArgs, file.Extension)
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Start()

	if !args.Persist {
		os.Exit(0)
	}
}

func getAllowedFiles(dir string) []File {
	var files []File

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && isAllowedExt(path) {
			stat, _ := os.Stat(path)

			files = append(files, File{
				Path:      path,
				Name:      stat.Name(),
				Extension: filepath.Ext(path),
			})
		}

		return nil
	})

	if err != nil {
		os.Exit(1)
	}

	return files
}

func isAllowedExt(file string) bool {
	allowed := []string{".png", ".jpg", ".jpeg"}

	ext := strings.ToLower(filepath.Ext(file))

	return slices.Contains(allowed, ext)
}

func exitAfterDelay(delay int64) {
	time.Sleep(time.Second * time.Duration(delay))

	os.Exit(0)
}

/*
* ClickableImage
 */

type ClickableImage struct {
	widget.BaseWidget

	image   *canvas.Image
	OnClick func()
}

func NewClickableImage(image image.Image) *ClickableImage {
	canvasImage := canvas.NewImageFromImage(image)

	ci := &ClickableImage{image: canvasImage}
	ci.image.ScaleMode = canvas.ImageScaleFastest
	ci.image.FillMode = canvas.ImageFillContain
	ci.ExtendBaseWidget(ci)

	return ci
}

func (ci *ClickableImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewPadded(ci.image))
}

func (ci *ClickableImage) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (ci *ClickableImage) Tapped(*fyne.PointEvent) {
	ci.OnClick()
}

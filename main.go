package main

import (
	"image"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/alexflint/go-arg"

	"wallpaper_picker/lib"
	"wallpaper_picker/lib/containers"
	"wallpaper_picker/lib/widgets"
)

var args struct {
	Grid    bool   `arg:"-g, --grid" help:"Display wallpapers in a grid."`
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

	files := lib.GetAllowedFiles(args.Dir)

	if len(files) == 0 {
		nowall := widget.NewLabelWithStyle("No wallpapers in directory.", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

		myWindow.SetContent(container.NewVBox(nowall))
		myWindow.Resize(fyne.NewSize(400, 80))

		go func() {
			time.Sleep(time.Second * time.Duration(5))

			os.Exit(0)
		}()
	} else {

		startupLabel := widget.NewLabelWithStyle("Loading Wallpapers", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		pb := widget.NewProgressBar()

		myWindow.SetContent(container.NewStack(startupLabel, pb))
		myWindow.Resize(fyne.NewSize(600, 80))

		go loadMainContent(files, myWindow, pb)
	}

	myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}

func loadMainContent(files []lib.File, window fyne.Window, progress *widget.ProgressBar) {
	total := len(files)
	var wg sync.WaitGroup
	done := make(chan bool, total)
	wg.Add(total)

	var wpContent *fyne.Container

	if args.Grid {
		wpContent = container.New(layout.NewGridLayout(3))
	} else {
		wpContent = container.New(layout.NewVBoxLayout())
	}

	// increased scroll factor, library default is too slow
	sc := containers.NewScaledScroll(fyne.ScrollVerticalOnly, 2.3, wpContent)

	for _, file := range files {
		go loadImage(file, wpContent, &wg, done)
	}

	go func() {
		wg.Wait()
		close(done)

		// Give the ui enough time to update progress
		time.Sleep(time.Millisecond * 10)
		fyne.Do(func() {
			window.Resize(fyne.NewSize(600, 486))
			window.SetContent(sc)
		})
	}()

	complete := 0
	for range done {
		complete++

		fyne.Do(func() {
			progress.SetValue(float64(complete) / float64(total))
		})
	}
}

func loadImage(f lib.File, c *fyne.Container, wg *sync.WaitGroup, done chan<- bool) {
	defer wg.Done()

	var image image.Image

	if args.Grid {
		image = lib.LoadSquareScaleImage(f)
	} else {
		image = lib.LoadWideScaleImage(f)
	}

	ci := widgets.NewClickableImage(image, func() {
		lib.SetWallpaper(f, args.WithExt, args.Command)

		if !args.Persist {
			os.Exit(0)
		}
	})

	if args.Grid {
		ci.SetMinSize(fyne.NewSize(186, 186))
	} else {
		ci.SetMinSize(fyne.NewSize(576, 162))
	}

	ci.Refresh()
	c.Add(ci)

	done <- true
}

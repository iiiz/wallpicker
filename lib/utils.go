package lib

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

	"golang.org/x/image/draw"
)

type File struct {
	Path      string
	Name      string
	Extension string
}

func GetAllowedFiles(dir string) []File {
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

// fixed to 32:9 scale
func LoadWideScaleImage(f File) image.Image {
	file, _ := os.Open(f.Path)
	defer file.Close()

	var original image.Image

	switch f.Extension {
	case ".jpeg", ".jpg":
		original, _ = jpeg.Decode(file)
	case ".png":
		original, _ = png.Decode(file)
	}

	originalBounds := original.Bounds()
	widthProduct := originalBounds.Dx() * 9
	heightProduct := originalBounds.Dy() * 32

	var drawingBounds image.Rectangle

	if widthProduct == heightProduct {
		// exact scale
		drawingBounds = originalBounds
	} else if widthProduct > heightProduct {
		// wider than 32:9, constrain by height
		newWidth := heightProduct / 9
		newMinX := (originalBounds.Max.X - newWidth) / 2
		newMaxX := newMinX + newWidth

		drawingBounds = image.Rectangle{
			Min: image.Point{
				X: newMinX,
				Y: originalBounds.Min.Y,
			},
			Max: image.Point{
				X: newMaxX,
				Y: originalBounds.Max.Y,
			},
		}
	} else {
		// taller than 32:9, constrain by width
		newHeight := widthProduct / 32
		newMinY := (originalBounds.Max.Y - newHeight) / 2
		newMaxY := newMinY + newHeight

		drawingBounds = image.Rectangle{
			Min: image.Point{
				X: originalBounds.Min.X,
				Y: newMinY,
			},
			Max: image.Point{
				X: originalBounds.Max.X,
				Y: newMaxY,
			},
		}
	}

	scaled := image.NewRGBA(image.Rect(0, 0, 960, 270))

	draw.NearestNeighbor.Scale(scaled, scaled.Rect, original, drawingBounds, draw.Over, nil)

	return scaled
}

func LoadSquareScaleImage(f File) image.Image {
	file, _ := os.Open(f.Path)
	defer file.Close()

	var original image.Image

	switch f.Extension {
	case ".jpeg", ".jpg":
		original, _ = jpeg.Decode(file)
	case ".png":
		original, _ = png.Decode(file)
	}

	bounds := original.Bounds()

	scaled := image.NewRGBA(image.Rect(0, 0, bounds.Max.Y/8, bounds.Max.Y/8))

	newMinX := (bounds.Max.X - bounds.Max.Y) / 2
	newMaxX := newMinX + bounds.Max.Y

	drawingBounds := image.Rectangle{
		Min: image.Point{
			X: newMinX,
			Y: bounds.Min.Y,
		},
		Max: image.Point{
			X: newMaxX,
			Y: bounds.Max.Y,
		},
	}

	draw.NearestNeighbor.Scale(scaled, scaled.Rect, original, drawingBounds, draw.Over, nil)

	return scaled
}

func SetWallpaper(file File, withExtension bool, command string) {
	commandFields := strings.Fields(command)
	cmdName := commandFields[0]
	cmdArgs := commandFields[1:]
	cmdArgs = append(cmdArgs, file.Path)

	if withExtension {
		cmdArgs = append(cmdArgs, file.Extension)
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Start()
}

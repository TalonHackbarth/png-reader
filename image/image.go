package image

import (
	goimage "image"
	"image/color"
)

type BaseImage interface {
	Width() int
	Height() int
	At(x int, y int) []int
	Clone() Image
	Convert() goimage.Image
}

// All Images Order will be [y][x][color]

const (
	chBinary    = 1
	chGrayScale = 1
	chPalette   = 1
	chRGB       = 3
	chRGBA      = 4
)

type Image struct {
	Data     [][]uint8
	Channels int
}

func (img *Image) Width() int {
	return len(img.Data[0]) / img.Channels
}

func (img *Image) Height() int {
	return len(img.Data)
}

func (img *Image) At(x int, y int) []int {
	result := make([]int, img.Channels)
	if (x*img.Channels)+img.Channels > len(img.Data[y]) {
		panic("Img Width not consistent with img channels")
	}
	for i := 0; i < img.Channels; i++ {
		result[i] = int(img.Data[y][(x*img.Channels)+i])
	}
	return result
}

func (img *Image) Convert() goimage.Image { // Only does Gray rn.
	output := goimage.NewGray(goimage.Rect(0, 0, img.Width(), img.Height()))
	output.Stride = img.Channels * 8
	for y := 0; y < img.Height(); y++ {
		for x := 0; x < img.Width(); x++ {
			output.SetGray(x, y, color.Gray{Y: uint8(img.At(x, y)[0])})
		}
	}
	return output
}

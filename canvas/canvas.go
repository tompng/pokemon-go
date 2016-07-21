package canvas

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"strings"
)

type Image interface {
	Get(x, y float64) (gray, alpha float64)
}

type ImageBuffer struct {
	Width  int
	Height int
	Gray   [][]float64
	Alpha  [][]float64
}

type SubImage struct {
	Source     Image
	X, Y, W, H float64
}

func (image *SubImage) Get(x, y float64) (float64, float64) {
	return image.Source.Get(image.X+image.W*x, image.Y+image.H*y)
}

func NewImageBufferFromFile(fileName string) *ImageBuffer {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Print("hoge")
		return nil
	}
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Print("piyo")
		return nil
	}
	rect := img.Bounds()
	image := NewImageBuffer(rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y)
	for x := 0; x < image.Width; x++ {
		for y := 0; y < image.Height; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			alpha := float64(a) / 0xffff
			gray := float64(r+g+b) / 3 / 0xffff
			image.Alpha[y][x] = alpha
			if alpha > 0 {
				image.Gray[y][x] = gray / alpha
			}
		}
	}
	return image
}

func NewImageBuffer(width int, height int) *ImageBuffer {
	gray := make([][]float64, height)
	alpha := make([][]float64, height)
	for y := 0; y < height; y++ {
		gray[y] = make([]float64, width)
		alpha[y] = make([]float64, width)
	}
	return &ImageBuffer{width, height, gray, alpha}
}

func (image *ImageBuffer) Get(x, y float64) (float64, float64) {
	if x < 0 || y < 0 || x > 1 || y > 1 {
		return 0, 0
	}
	ix := int(x * float64(image.Width))
	if ix >= image.Width {
		ix = image.Width - 1
	}
	iy := int(y * float64(image.Height))
	if iy >= image.Height {
		iy = image.Height - 1
	}
	return image.Gray[iy][ix], image.Alpha[iy][ix]
}

func (image *ImageBuffer) Sub(x, y, w, h float64) *SubImage {
	return &SubImage{image, x, y, w, h}
}

func (image *ImageBuffer) Plot(x, y int, gray, alpha float64) {
	if x < 0 || y < 0 || x >= image.Width || y >= image.Height {
		return
	}
	dstGray, dstAlpha := image.Gray[y][x], image.Alpha[y][x]
	newAlpha := dstAlpha + alpha - dstAlpha*alpha
	image.Alpha[y][x] = newAlpha
	if newAlpha == 0 {
		image.Gray[y][x] = 0
	} else {
		image.Gray[y][x] = (dstGray*dstAlpha*(1-alpha) + gray*alpha) / newAlpha
	}
}

func (screen *ImageBuffer) Draw(image Image, x, y, w, h float64) {
	if x+w < 0 || y+h < 0 || float64(screen.Width) < x || float64(screen.Height) < y {
		return
	}
	for ix := int(x); float64(ix) < x+w; ix++ {
		for iy := int(y); float64(iy) < y+h; iy++ {
			gray, alpha := image.Get((float64(ix)-x)/w, (float64(iy)-y)/h)
			screen.Plot(ix, iy, gray, alpha)
		}
	}
}

func (image *ImageBuffer) String() string {
	lines := make([]string, image.Height/2)
	for y := 0; y < image.Height/2; y++ {
		var buf bytes.Buffer
		for x := 0; x < image.Width; x++ {
			ug, ua := image.Gray[2*y][x], image.Alpha[2*y][x]
			up := int(16 * (ug*ua + 1*(1-ua)))
			dg, da := image.Gray[2*y+1][x], image.Alpha[2*y+1][x]
			down := int(16 * (dg*da + 1*(1-da)))
			if up < 0 {
				up = 0
			}
			if down < 0 {
				down = 0
			}
			if up > 0xf {
				up = 0xf
			}
			if down > 0xf {
				down = 0xf
			}
			buf.WriteByte(charTable[up][down])
		}
		lines[y] = buf.String()
	}
	return strings.Join(lines, "\n")
}

func (image *ImageBuffer) Print() {
	fmt.Print("\x1B[1;1H", image.String())
}

var charTable []string = []string{
	"MMMMMM###TTTTTTT",
	"QQBMMNW##TTTTTV*",
	"QQQBBEK@PTTTVVV*",
	"QQQmdE88P9VVVV**",
	"QQQmdGDU0YVV77**",
	"pQQmAbk65YY?7***",
	"ppgAww443vv?7***",
	"pggyysxcJv??7***",
	"pggyaLojrt<<+**\"",
	"gggaauuj{11!//\"\"",
	"gggaauui])|!/~~\"",
	"ggaauui]((;::~~^",
	"ggaauu](;;::-~~'",
	"ggauu(;;;;---~``",
	"gaau;;,,,,,...``",
	"gau,,,,,,,,...  "}

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tty "github.com/mattn/go-tty"
	"github.com/tompng/pokemon-go"
	"github.com/tompng/pokemon-go/canvas"
)

var fontData *canvas.ImageBuffer

func DrawString(screen *canvas.ImageBuffer, message string, x, y, size float64) {
	if fontData == nil {
		fontData = canvas.NewImageBufferFromReader(file("images/chars.png"))
	}
	for i, c := range message {
		face := fontData.Sub(float64(c%16)/16.0, float64(c/16)/8.0, 1/16.0, 1/8.0)
		screen.Draw(face, x+float64(i)*size/2, y, size/2, size)
	}
}

func FileNames(path string) []string {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}
	}
	files := make([]string, len(fileInfos))
	for i, fi := range fileInfos {
		file := fi.Name()
		if fi.IsDir() {
			file += "/"
		}
		files[i] = file
	}
	return files
}

func PokemonImage() *canvas.ImageBuffer {
	currentPath, err := filepath.Abs(".")
	if err != nil {
		panic("cannot get current path")
	}
	var seed int64
	for i, c := range currentPath {
		seed = seed*0x987654321 + int64(c) + int64(i)
	}
	var imageFiles []string
	for _, f := range pokemongo.AssetNames() {
		if strings.HasPrefix(f, "images/pokemon/") && path.Ext(f) == ".png" {
			imageFiles = append(imageFiles, f)
		}
	}
	sort.Strings(imageFiles)
	random := rand.New(rand.NewSource(seed))
	f := imageFiles[random.Intn(len(imageFiles))]
	return canvas.NewImageBufferFromReader(file(f))
}

func DrawScrollingMessages(screen *canvas.ImageBuffer, messages []string, time float64) {
	for i, message := range messages {
		DrawString(screen, message, 0, float64(i+1)*16-20*time, 16)
	}
}

func DrawGotcha(screen *canvas.ImageBuffer) {
	message := "Gotcha!"
	width := 80
	length := len(message)
	height := 2 * width / length
	r := 4.0
	for x := -int(r); x < width+int(r); x++ {
		for y := -int(r); y < height; y++ {
			dx, dy := 1.0, 1.0
			if x < 0 {
				dx = 1.0 + float64(x)/r
			}
			if x > width {
				dx = 1.0 - float64(x-width)/r
			}
			if y < 0 {
				dy = 1.0 + float64(y)/r
			}
			screen.Plot(screen.Width/2-width/2+x, screen.Height-height+y, 1, 0.8*dx*dy)
		}
	}
	DrawString(screen, message, float64(screen.Width-width)/2, float64(screen.Height-height), float64(height))
}

func Save(pokemon canvas.Image) {
	screen := canvas.NewImageBuffer(80, 80)
	screen.Draw(pokemon, 0, 0, 80, 80)
	DrawGotcha(screen)
	ioutil.WriteFile("pokemon.txt", []byte(screen.String()), 0666)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%v: %v", os.Args[0], err)
	os.Exit(1)
}

func file(n string) io.Reader {
	b, err := pokemongo.Asset(n)
	if err != nil {
		panic(err.Error())
	}
	return bytes.NewReader(b)
}

func main() {
	tty, err := tty.Open()
	if err != nil {
		fatal(err)
	}

	ball1 := canvas.NewImageBufferFromReader(file("images/ball1.png"))
	smoke := canvas.NewImageBufferFromReader(file("images/smoke.png"))
	ball2 := canvas.NewImageBufferFromReader(file("images/ball2.png"))
	ball3 := canvas.NewImageBufferFromReader(file("images/ball3.png"))

	fileNames := FileNames(".")
	pokemon := PokemonImage()

	time0 := time.Now().UnixNano()
	rand.Seed(time0)

	dstx, dsty := 0.5, 0.5
	dx1, dy1 := 2*rand.Float64()-1, 2*rand.Float64()-1
	dx2, dy2 := 2*rand.Float64()-1, 2*rand.Float64()-1
	getTime := 20.0
	exitFlag := false
	currentTime := func() float64 {
		return float64(time.Now().UnixNano()-time0) / 1000 / 1000 / 1000
	}

	go func() {
		for {
			code, _ := tty.ReadRune()
			if code == 0x03 || code == 0x1C {
				exitFlag = true
				continue
			}
			time := currentTime()
			if time < 2.0 {
				time = 2.0
			}
			if time < getTime {
				getTime = time
			}
		}
	}()
	throwTime := 1.0
	for {
		terminalWidth, terminalHeight, err := tty.Size()
		if err != nil {
			fatal(err)
		}
		t := currentTime()
		screen := canvas.NewImageBuffer(terminalWidth, (terminalHeight-1)*2)

		size := 96 * (1 - math.Exp(-t)) * (1 + 0.1*(math.Sin(1.4*t)+math.Sin(1.9*t)))
		x := dstx + dx1*math.Exp(-t) + dx2*math.Exp(-t/2) + 0.1*(math.Sin(2.1*t)+math.Sin(1.7*t))
		y := dsty + dy1*math.Exp(-t) + dy2*math.Exp(-t/2) + 0.1*(math.Sin(2.5*t)+math.Sin(1.3*t))
		if t < getTime+throwTime {
			screen.Draw(pokemon, float64(screen.Width)*x-size/2, float64(screen.Height)*y-size/2, size, size)
		}
		DrawScrollingMessages(screen, fileNames, t)

		if getTime < t {
			phase := t - getTime
			if phase < throwTime {
				x := size/4*(2*rand.Float64()-1) + float64(screen.Width)/2 - size/2
				y := size/4*(2*rand.Float64()-1) + float64(screen.Height)/2 - size/2
				screen.Draw(smoke, x, y, size, size)
			}
			yt := phase * 2
			if yt < 0.5 {
				yt = 0.5
			}
			ytmod1 := yt - math.Floor(yt)
			pos := 4 * ytmod1 * (1 - ytmod1) * math.Exp(-math.Floor(yt))
			if phase > 2 {
				pos = 0
			}
			ball := ball1
			if phase > 5 {
				ball = ball3
			} else if phase > 1 && phase-math.Floor(phase) < 0.1 {
				ball = ball2
			}
			screen.Draw(ball, float64(screen.Width)/2, (float64(screen.Height)/2+20)*(1-pos)-20, 20, 20)
			if phase > 5 {
				DrawGotcha(screen)
				screen.Print()
				Save(pokemon)
				break
			}
		}
		screen.Print()
		if exitFlag {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

}

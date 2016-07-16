package main

import (
	"time"
	"math"
	"math/rand"
	"./canvas"
)

var fontData *canvas.ImageBuffer
func DrawString(screen *canvas.ImageBuffer, message string, x, y, size float64){
  if fontData == nil{fontData = canvas.NewImageBufferFromFile("images/chars.png")}
  for i, c := range message {
    face := fontData.Sub(float64(c%16)/16.0, float64(c/16)/8.0, 1/16.0, 1/8.0)
    screen.Draw(face, x+float64(i)*size/2, y, size/2, size)
  }
}

func PokemonImage() *canvas.ImageBuffer{
  return canvas.NewImageBufferFromFile("images/pokemon/gopher.png")
}

func DrawScrollingMessages(screen *canvas.ImageBuffer, messages []string, time float64){
	for i, message := range messages {
		DrawString(screen, message, 0, float64(i+1)*16 - 20*time, 16)
	}
}

func DrawGotcha(screen *canvas.ImageBuffer){
	message := "Gotcha!"
	width := 80
	length := len(message)
	height := 2*width/length
	r := 4.0
	for x:= -int(r); x < width + int(r); x++ {
		for y:= -int(r); y<height; y++{
			dx, dy := 1.0, 1.0
			if x < 0 {dx = 1.0 + float64(x)/r}
			if x > width {dx = 1.0 - float64(x - width)/r}
			if y < 0 {dy = 1.0 + float64(y)/r}
			screen.Plot(screen.Width/2-width/2+x, screen.Height-height+y, 1, 0.8*dx*dy)
		}
	}
	DrawString(screen, message, float64(screen.Width-width)/2, float64(screen.Height-height), float64(height))
}

func main() {
  ball1 := canvas.NewImageBufferFromFile("images/ball1.png")
  smoke := canvas.NewImageBufferFromFile("images/smoke.png")
  ball2 := canvas.NewImageBufferFromFile("images/ball2.png")
  ball3 := canvas.NewImageBufferFromFile("images/ball3.png")
  pokemon := PokemonImage()

  messages := []string{"foo", "bar", "hoge", "piyo"}

	time0 := time.Now().UnixNano()
	rand.Seed(time0)

	dstx, dsty := 0.5, 0.5
	dx1, dy1 := 2*rand.Float64()-1, 2*rand.Float64()-1
	dx2, dy2 := 2*rand.Float64()-1, 2*rand.Float64()-1
	getTime := 4.0
	throwTime := 1.0
	for ;; {
		t := float64(time.Now().UnixNano() - time0)/1000/1000/1000
		screen := canvas.NewImageBuffer(80, 80)

		size := 96*(1-math.Exp(-t))*(1+0.1*(math.Sin(1.4*t)+math.Sin(1.9*t)))
		x := dstx+dx1*math.Exp(-t)+dx2*math.Exp(-t/2)+0.1*(math.Sin(2.1*t)+math.Sin(1.7*t))
		y := dsty+dy1*math.Exp(-t)+dy2*math.Exp(-t/2)+0.1*(math.Sin(2.5*t)+math.Sin(1.3*t))
		if t < getTime + throwTime {
			screen.Draw(pokemon, float64(screen.Width)*x-size/2, float64(screen.Height)*y-size/2, size, size)
		}
		DrawScrollingMessages(screen, messages, t)

		if getTime < t {
			phase := t - getTime
			if phase < throwTime {
				x := size/4*(2*rand.Float64()-1)+float64(screen.Width)/2-size/2
				y := size/4*(2*rand.Float64()-1)+float64(screen.Height)/2-size/2
				screen.Draw(smoke, x, y, size, size)
			}
			yt := phase*2
			if yt < 0.5 {
				yt = 0.5
			}
			ytmod1 := yt - math.Floor(yt)
			pos := 4*ytmod1*(1-ytmod1)*math.Exp(-math.Floor(yt))
			if phase > 2 {
				pos = 0
			}
			ball := ball1
			if phase > 5 {
				ball = ball3
			} else if phase > 1 && phase - math.Floor(phase) < 0.1 {
				ball = ball2
			}
			screen.Draw(ball, float64(screen.Width)/2,(float64(screen.Height)/2+20)*(1-pos)-20, 20, 20)
			if phase > 5 {
				DrawGotcha(screen)
			}
		}
		screen.Print()
		time.Sleep(50*time.Millisecond)
	}

}

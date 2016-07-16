package main

import (
	"fmt"
	"time"
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

func main() {
	image := canvas.NewImageBufferFromFile("images/ball1.png")
	fmt.Print(image.Width)
	// screen := canvas.NewImageBuffer(80, 80)
	// screen.Draw(image.Sub(0.5, 0.5, 0.5, 0.5), 0, 0, 80, 80)
	// screen.Draw(image, 30, 50, 40, 40)
  // DrawString(screen, "hello", 0, 0, 20)
	// screen.Print()


  // ball1 := canvas.NewImageBufferFromFile("images/ball1.png")
  // smoke := canvas.NewImageBufferFromFile("images/smoke.png")
  // ball2 := canvas.NewImageBufferFromFile("images/ball2.png")
  // ball3 := canvas.NewImageBufferFromFile("images/ball3.png")
  // pokemon := PokemonImage()
	//
  messages := []string{"foo", "bar", "hoge", "piyo"}

	time0 := time.Now().UnixNano()
	for ;; {
		t := float64(time.Now().UnixNano() - time0)/1000/1000/1000
		screen := canvas.NewImageBuffer(80, 80)


		DrawScrollingMessages(screen, messages, t)
		screen.Print()
		time.Sleep(50*time.Millisecond)
	}
	fmt.Print(time0)

}

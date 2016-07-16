package main

import (
	"./canvas"
	"fmt"
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


func main() {
	image := canvas.NewImageBufferFromFile("images/ball1.png")
	fmt.Print(image.Width)
	screen := canvas.NewImageBuffer(80, 80)
	// screen.Draw(image.Sub(0.5, 0.5, 0.5, 0.5), 0, 0, 80, 80)
	screen.Draw(image, 30, 50, 40, 40)
  DrawString(screen, "hello", 0, 0, 20)
	screen.Print()

  // ball1 := canvas.NewImageBufferFromFile("images/ball1.png")
  // smoke := canvas.NewImageBufferFromFile("images/smoke.png")
  // ball2 := canvas.NewImageBufferFromFile("images/ball2.png")
  // ball3 := canvas.NewImageBufferFromFile("images/ball3.png")
  // pokemon := PokemonImage()
	//
  // messages := []string{"foo", "bar", "hoge", "piyo"}


}

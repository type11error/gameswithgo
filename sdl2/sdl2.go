package main

import (
  "fmt"
  "github.com/veandco/go-sdl2/sdl"
)

const winWidth int = 800
const winHeight int = 600

type color struct {
  r, g, b byte
}

func setPixel(x, y int, c color, pixels []byte) {
  index := (y*winWidth + x) * 4

  if index < len(pixels)-4 && index >= 0 {
    pixels[index] = c.r
    pixels[index+1] = c.g
    pixels[index+2] = c.b
  }
}


func main() {

  fmt.Println("Starting...")

  err := sdl.Init(sdl.INIT_EVERYTHING)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer sdl.Quit()

  window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
    int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)

  if err != nil {
    fmt.Println(err)
    return
  }
  defer window.Destroy()

  renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer renderer.Destroy()

  tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888,sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
  if err != nil {
    fmt.Println(err)
    return
  }
  defer tex.Destroy()

  pixels := make([]byte, winWidth*winHeight*4)

  for y := 0; y < winHeight; y++ {
    for x := 0; x < winWidth; x++ {
      setPixel(x,y, color{255,0,0}, pixels)
    }
  }

  tex.Update(nil, pixels, winWidth*4)
  renderer.Copy(tex, nil, nil)
  renderer.Present()


  running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}

}

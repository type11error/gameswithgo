package main

import (
  "fmt"
  "github.com/veandco/go-sdl2/sdl"
  "image/png"
  "os"
  "time"
  "github.com/type11error/gameswithgo/noise"
)

const winWidth int = 800
const winHeight int = 600

type balloon struct {
  tex *sdl.Texture
  pos
  scale float32
  w,h int
}

func(balloon *balloon) draw(renderer *sdl.Renderer) {
  newW := int32(float32(balloon.w)*scale)
  newH := int32(float32(balloon.h)*scale)

  x := int32(balloon.x - float32(newW)/2)
  y := int32(balloon.y - float32(newH)/2)
  rect := &sdl.Rect{x,y,newW, newH}
  renderer.Copy(baloon.tex, nil rect)
}

type rgba struct {
  r, g, b byte
}

type pos struct {
  x, y float32
}

func setPixel(x, y int, c rgba, pixels []byte) {
  index := (y*winWidth + x) * 4

  if index < len(pixels)-4 && index >= 0 {
    pixels[index] = c.r
    pixels[index+1] = c.g
    pixels[index+2] = c.b
  }
}



func clear(pixels []byte) {
  for i := range pixels {
    pixels[i] = 0
  }
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte,w,h,int) *sdl.Texture {
  tex,err := renderer.CreateTexture(sd.PIXELFORMAT_ABGR8888,sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))

  if err != nil {
    panic(err)
  }

  tex.Update(nil, pixels, w*4)
  return tex

}

func loadBalloons(renderer *sdl.Renderer) []baloons{

  balloonStrs := []string{"balloon_red.png", "balloon_green.png", "balloon_blue.png"}
  balloons := make([]balloons, len(balloonStrs))

  for i, bstr := range balloonStrs {

    infile,err := os.Open(bstr)
    if err != nil {
      panic(err)
    }

    img,err := png.Decode(infile)
    if err != nil {
      panic(err)
    }

    w := img.Bounds().Max.X
    h := img.Bounds().Max.Y

    balloonPixels := make([]byte,w*h*4)
    bIndex := 0
    for y := 0; y<h; y++ {
      for x := 0; x < w; x++ {
        r,g,b,a := img.At(x,y).RGBA()
        balloonPixels[bIndex] = byte(r/256)
        bIndex++
        balloonPixels[bIndex] = byte(g/256)
        bIndex++
        balloonPixels[bIndex] = byte(b/256)
        bIndex++
        balloonPixels[bIndex] = byte(a/256)
        bIndex++
      }
    }

    tex := pixelsToTexture(renderer, balloonPixels,w,h)
    balloons[i] = balloon{tex, pos{float32(i*60),float32(i*60)}, float32(1+i), w, h}
  }

  return balloons
}

// calculate a value between two numbers given a percentage
func flerp(b1 byte, b2 byte, pct float32) byte {
  return byte(float32(b1) + pct*(float32(b2)-float32(b1)))
}


// lerp for a rgba value
func rgbaLerp(c1, c2 rgba, pct float32) rgba{
  return rgba{flerp(c1.r, c2.r, pct), flerp(c1.g, c2.g, pct), flerp(c1.b, c2.b, pct)}
}

// calculate a 256 gradient between two rgbas
func getGradient(c1, c2 rgba) []rgba {
  result := make([]rgba, 256)
  for i := range result {
    pct := float32(i) / float32(255)
    result[i] = rgbaLerp(c1, c2, pct)
  }
  return result
}

func getDualGradient(c1,c2,c3,c4 rgba) [] rgba {
  result := make([]rgba, 256)
  for i := range result {
    pct := float32(i) / float32(255)
    if pct < 0.5 {
      result[i] = rgbaLerp(c1, c2, pct*float32(2))
    } else {
      result[i] = rgbaLerp(c3, c4, pct*float32(1.5)-float32(0.5))
    }
  }
  return result
}

// make a value only between a min and max
func clamp(min, max, v int) int {
  if v < min {
    v = min
  } else if v > max {
    v = max
  }

  return v
}

func rescaleAndDraw(noise []float32, min, max float32, gradient []rgba, w, h int) []byte {
  result := make([]byte,w*h*4)
  scale := 255.0 / (max - min)
  offset := min * scale

  for i := range noise {
    noise[i] = noise[i]*scale - offset
    c := gradient[clamp(0, 255, int(noise[i]))]
    p := i * 4
    result[p] = c.r
    result[p+1] = c.g
    result[p+2] = c.b
  }

  return result
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

  cloudNoise, min, max := noise.MakeNoise(noise.FBM, .009, .5, 3, 3, winWidth, winHeight)
  cloudGradient := getGradient(rgba{0,0,255}, rgba{255,255,255})
  cloudPixels := rescaleAndDraw(cloudNoise, min, max, cloudGradient, winWidth, winHeight)
  cloudTexture := pixelsToTexture(renderer, cloudPixels, winWidth, winHeight)

  pixels := make([]byte, winWidth*winHeight*4)
  balloonTextures := loadBalloons(renderer)



  running := true

  dir := 1
  for running {
    frameStart := time.Now()

    for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
      switch event.(type) {
      case *sdl.QuitEvent:
        println("Quit")
        running = false
        break
      }
    }

    renderer.Copy(cloudTexture, nil, nil)

    for _, baloon := range balloons {
      balloon.draw()
    }

    balloonTextures[1].x += float32(1 *dir)
    if balloonTextures[1].x > 400 || balloonTextures[1].x < 0 {
      dir = dir * -1
    }

    tex.Update(nil, pixels, winWidth*4)
    renderer.Copy(tex, nil, nil)
    renderer.Present()

    elapsedTime := float32(time.Since(frameStart).Seconds()*1000)
    fmt.Println("ms per frame:", elapsedTime) 
    if elapsedTime < 5 {
      sdl.Delay(4 - uint32(elapsedTime))
      elapsedTime = float32(time.Since(frameStart).Seconds())
    }
  }

}

package main

import "fmt"

type position struct {
  x float32
  y float32
}

type badGuy struct {
  name string
  health int
  pos position
}


func addOne(num *int) {

  *num = *num + 1

}

func main() {

  x := 5
  fmt.Println(x)

  var xPtr *int = &x
  fmt.Println(xPtr)

  addOne(xPtr)
  fmt.Println(x)



}

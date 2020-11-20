package main

// tries it took 
// tell if the user is lying
//

import (
  "fmt"
  "bufio"
  "os"
)

func main() {

  scanner := bufio.NewScanner(os.Stdin)

  fmt.Println("Please think of a number between 1 and 100")
  fmt.Println("Press ENTER when ready")
  scanner.Scan()

  low := 1
  high := 100
  tries := 0

  for {
    guess := (low + high) / 2
    fmt.Println("I guess the number is", guess)
    fmt.Println("Is that:")
    fmt.Println("(a) too high?")
    fmt.Println("(b) too low?")
    fmt.Println("(c) correct?")
    scanner.Scan()
    response := scanner.Text()

    if response == "a" {
      high = guess - 1
      tries++
    } else if response == "b" {
      low = guess + 1
      tries++
    } else if response == "c" {
      fmt.Println("I won with tries: ", tries)
      break
    } else {
      fmt.Println("Invalid response, try again.")
    }

    if low > high {
      fmt.Println("We guessed your number you are lying")
      break
    }


  }

}

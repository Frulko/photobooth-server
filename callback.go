package main

import "fmt"

func main() {
  makeFibGen(handleCallback)
}

func handleCallback(i int) {
  fmt.Println("callback called", i)
} 


func makeFibGen(callback func(i int)) {
  f1 := 0
  f2 := 1
 
  for i := 0; i < 10; i++ {
    f2, f1 = (f1 + f2), f2
    callback(f1)
  }
}
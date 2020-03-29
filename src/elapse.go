package main

import (
  "time"
  "fmt"
)

func elapsed(what string) func() {
    start := time.Now()
    fmt.Println("start")
    return func() {
        fmt.Printf("%s took %v\n", what, time.Since(start))
    }
}

func main() {
    defer elapsed("page")()
    time.Sleep(time.Second * 3)
}

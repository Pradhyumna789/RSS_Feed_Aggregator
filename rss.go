package main

import (
  "fmt"
  "github.com/Pradhyumna789/RSS_Aggregator/internal/config"
)

func main() {
  var configStruct config.Config

  initial := config.Read()
  configStruct.SetUser("pradyumna")
  updated := config.Read()

  fmt.Println(initial)
  fmt.Println(updated)
}


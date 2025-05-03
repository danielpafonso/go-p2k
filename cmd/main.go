package main

import (
	"fmt"

	"go-p2k/internal"
)

func main() {
	fmt.Println("Start")

	configs, err := internal.LoadConfigurations("./config.json")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", configs)

	fmt.Println("End")
}

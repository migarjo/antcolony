package main

import (
	"fmt"
)

type model struct{}

func main() {
	home := createTowns(1, 100)
	towns := createTowns(10, 100)

	fmt.Println(home)
	fmt.Println(towns)
}

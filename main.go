package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!!!")

	z := make([]*int, 20)
	x := 125
	z[0] = &x
	println(z)

	// stdin := os.Stdin

}

package main

import (
	"fmt"
	"strings"
)

func main() {

	var edad int

	fmt.Println("Introduce una edad")

	fmt.Scan(&edad)

	edadBin := fmt.Sprintf("%b\n", edad)

	fmt.Println(strings.Count(edadBin, "1"))
}

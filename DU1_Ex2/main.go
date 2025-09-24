package main

import (
	"fmt"
	"strings"
)

func main() {
	var numeros string

	fmt.Scanln(&numeros)

	numerosSueltos := strings.Split(numeros, " ")

	fmt.Println(len(numerosSueltos))

	fmt.Println(numerosSueltos)

}

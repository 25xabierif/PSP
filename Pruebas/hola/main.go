package main

import "fmt"

func main() {

	continuar := true

	for continuar {
		fmt.Println("Hola, Mundo!")
		continuar = false
	}

	for i := 0; i <= 10; i++ {
		fmt.Println("Hola, Mundo!")
	}

	//map con make
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3

	fmt.Println(dereference(4))

}

func dereference(x int) int {
	x = x ^ 2
	return x
}

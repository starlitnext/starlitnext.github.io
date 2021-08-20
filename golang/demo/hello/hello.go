package main

import (
	"fmt"

	"example.com/grettings"
)

func main() {
	message := grettings.Hello("Gladys")
	fmt.Println(message)
}

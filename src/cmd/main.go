package main

import (
	"fmt"
	"game-message-core/jsonData"
)

func main() {
	fmt.Printf("this is test main go \n")
	v := jsonData.Vector3{X: 1, Z: 999.99}
	fmt.Println(v)
}

package main

import (
	"fmt"
	"game-message-core/jsonData"

	"github.com/Meland-Inc/game-services/src/global/configModule"
)

func main() {
	fmt.Printf("this is test main go \n")
	v := jsonData.Vector3{X: 1, Z: 999.99}
	fmt.Println(v)

	err := configModule.Init()
	fmt.Println(err)
}

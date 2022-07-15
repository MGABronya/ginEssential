package main

import (
	"fmt"
	"os"
)

func main() {
	dir, _ := os.Getwd()
	fmt.Print(dir)

}

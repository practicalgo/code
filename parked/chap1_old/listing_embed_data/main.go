package main

import "fmt"
import "github.com/practicalgo/code/chap1/listing7/assets"

func main() {
	htmlData := assets.FileData["./assets/index.html"]
	fmt.Printf("%s", string(htmlData))
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	code := fmt.Sprintf(`
// ^// Code generated .* DO NOT EDIT\.$
package assets
var FileData map[string][]byte = map[string][]byte{
    "%s": %#v,
}
`, os.Args[1], data)
	err = ioutil.WriteFile("./assets/assets.go", []byte(code), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

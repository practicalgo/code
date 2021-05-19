package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	u, err := url.Parse("http://user:pass@example.com/api/?name=\"jane doe\"&age=25&name=john#page1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", u.Path)
	fmt.Printf("%#v\n", u.Query())
	fmt.Printf("%#v\n", u.Path)
	fmt.Printf("%#v\n", u.User.Username())
	p, ok := u.User.Password()
	fmt.Printf("%v %v\n", p, ok)

}

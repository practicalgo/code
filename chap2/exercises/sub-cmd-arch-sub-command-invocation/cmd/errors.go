// Listing 1.13: chap1/listing5/cmd/errors.go
package cmd

import "errors"

var ErrNoServerSpecified = errors.New("You have to specify the remote server.")

package cmd

import "errors"

var ErrNoServerSpecified = errors.New("you have to specify the remote server")
var ErrInvalidRegisterArguments = errors.New("you have to specify the package name, version and file path to upload")
var ErrInvalidQueryArguments = errors.New("you have to specify the package name to query")

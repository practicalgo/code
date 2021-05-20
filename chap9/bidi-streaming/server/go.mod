module github.com/practicalgo/code/chap9/bidi-streaming/server

go 1.16

require google.golang.org/grpc v1.37.0 // indirect

require github.com/practicalgo/code/chap9/bidi-streaming/service v0.0.0

replace github.com/practicalgo/code/chap9/bidi-streaming/service => ../service

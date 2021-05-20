module github.com/practicalgo/code/chap9/bindata-client-streaming/server

go 1.16

require google.golang.org/grpc v1.37.0

require github.com/practicalgo/code/chap9/bindata-client-streaming/service v0.0.0

replace github.com/practicalgo/code/chap9/bindata-client-streaming/service v0.0.0 => ../service

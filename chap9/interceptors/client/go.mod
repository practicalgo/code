module github.com/practicalgo/code/chap9/interceptors/client

go 1.16

require (
	github.com/practicalgo/code/chap9/interceptors/service v0.0.0
	google.golang.org/grpc v1.37.0
)

replace github.com/practicalgo/code/chap9/interceptors/service => ../service

module github.com/practicalgo/code/chap10/panic-svc-handling/client

go 1.16

require (
	github.com/practicalgo/code/chap10/panic-svc-handling/service v0.0.0
	google.golang.org/grpc v1.37.0
)

replace github.com/practicalgo/code/chap10/panic-svc-handling/service => ../service

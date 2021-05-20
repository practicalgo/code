module github.com/practicalgo/code/chap8/user-sevice/client

go 1.16

require (
	github.com/practicalgo/code/chap8/user-service/service-v1 v0.0.0
	google.golang.org/grpc v1.37.0 // indirect
)

replace github.com/practicalgo/code/chap8/user-service/service-v1 => ../service-v1

module github.com/practicalgo/code/chap8/exercises/mync-grpc

require (
	github.com/practicalgo/code/chap8/user-service/service v0.0.0
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.25.0 // indirect
)

replace github.com/practicalgo/code/chap8/user-service/service => ../../user-service/service

go 1.16

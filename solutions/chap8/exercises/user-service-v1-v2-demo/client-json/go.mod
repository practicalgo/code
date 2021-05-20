module github.com/practicalgo/code/chap8/user-sevice/client-json

go 1.16

require (
	github.com/practicalgo/code/chap8/user-service/service v0.0.0
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)

replace github.com/practicalgo/code/chap8/user-service/service => ../service

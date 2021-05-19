module github.com/practicalgo/code/chap8/user-sevice/server

go 1.16

require google.golang.org/grpc v1.37.1

require github.com/practicalgo/code/chap8/user-service/service v0.0.0

replace github.com/practicalgo/code/chap8/user-service/service => ../service

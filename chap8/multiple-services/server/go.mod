module github.com/practicalgo/code/chap8/multiple-sevices/server

go 1.16

require google.golang.org/grpc v1.37.1

require github.com/practicalgo/code/chap8/multiple-services/service v0.0.0

replace github.com/practicalgo/code/chap8/multiple-services/service => ../service

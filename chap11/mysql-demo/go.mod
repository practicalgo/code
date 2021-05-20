module github.com/practicalgo/code/chap11/mysql-demo

go 1.16

require (
	github.com/docker/go-connections v0.4.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/testcontainers/testcontainers-go v0.11.1
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/net v0.0.0-20210505214959-0714010a04ed // indirect
	golang.org/x/sys v0.0.0-20210503173754-0981d6026fa6 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	google.golang.org/genproto v0.0.0-20210506142907-4a47615972c2 // indirect
	google.golang.org/grpc v1.37.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

// Remove replace and upgrade library once
// https://github.com/testcontainers/testcontainers-go/pull/342 is merged
// The tag used here is on my personal fork containing the change in PR:
// https://github.com/amitsaha/testcontainers-go/releases/tag/v0.11.1-pr-342
replace github.com/testcontainers/testcontainers-go => github.com/amitsaha/testcontainers-go v0.11.1-pr-342

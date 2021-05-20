module github.com/practicalgo/code/chap11/pkg-server

go 1.16

require (
	github.com/docker/go-connections v0.4.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/testcontainers/testcontainers-go v0.11.1
	gocloud.dev v0.23.0
)

// Remove replace and upgrade library once
// https://github.com/testcontainers/testcontainers-go/pull/342 is merged
// The tag used here is on my personal fork containing the change in PR:
// https://github.com/amitsaha/testcontainers-go/releases/tag/v0.11.1-pr-342
replace github.com/testcontainers/testcontainers-go => github.com/amitsaha/testcontainers-go v0.11.1-pr-342

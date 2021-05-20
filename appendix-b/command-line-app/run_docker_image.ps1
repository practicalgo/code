docker run -v ${PWD}:/data -e X_AUTH_TOKEN=token-123 `
        --network appendix-b_default `
        -ti practicalgo/pkgcli `
        register `
        -name "test" -version 0.1 -path /data/main.go `
        http://pkgserver:8080
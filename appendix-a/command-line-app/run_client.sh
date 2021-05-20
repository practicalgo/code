X_AUTH_TOKEN=token-123 ./pkgcli -log-level=0 -jaeger-addr http://127.0.0.1:14268 -metrics-addr 127.0.0.1:9125  register -name "test" -version 0.7 -path main.go http://127.0.0.1:8080 


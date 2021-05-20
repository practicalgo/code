$env:X_AUTH_TOKEN="token-123"
.\pkgcli.exe register -name "test" -version 0.7 -path main.go http://127.0.0.1:8080


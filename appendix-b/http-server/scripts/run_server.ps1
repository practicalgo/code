$env:LOG_LEVEL=0
$env:JAEGER_ADDR="http://127.0.0.1:14268"
$env:STATSD_ADDR="127.0.0.1:9125"
$env:AWS_ACCESS_KEY_ID="admin"
$env:AWS_SECRET_ACCESS_KEY="admin123"
$env:BUCKET_NAME="test-bucket"
$env:S3_ADDR="localhost:9000"
$env:DB_ADDR="localhost:3306"
$env:DB_NAME="package_server"
$env:DB_USER="packages_rw"
$env:DB_PASSWORD="password"
$env:USERS_SVC_ADDR="localhost:50051"

.\http-server.exe
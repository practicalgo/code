LOG_LEVEL=0 \
JAEGER_ADDR=http://127.0.0.1:14268 STATSD_ADDR=127.0.0.1:9125 \
AWS_ACCESS_KEY_ID=admin AWS_SECRET_ACCESS_KEY=admin123 \
BUCKET_NAME=test-bucket S3_ADDR=localhost:9000 \
DB_ADDR=localhost:3306 DB_NAME=package_server \
DB_USER=packages_rw DB_PASSWORD=password USERS_SVC_ADDR=localhost:50051 ./http-server
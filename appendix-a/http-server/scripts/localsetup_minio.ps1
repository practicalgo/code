docker run `
    -p 9000:9000 `
    -p 9001:9001 `
    -e MINIO_ROOT_USER=admin `
    -e MINIO_ROOT_PASSWORD=admin123 `
    -ti minio/minio:RELEASE.2021-07-08T01-15-01Z `
    server "/data" `
    --console-address ":9001"

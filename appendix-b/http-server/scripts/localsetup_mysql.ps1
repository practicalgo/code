$BOOTSTRAP_SQL_PATH="${PWD}\mysql-init"

docker run `
    --platform linux/x86_64 `
    -p 3306:3306 `
    -e MYSQL_ROOT_PASSWORD=rootpassword `
    -e MYSQL_DATABASE=package_server `
    -e MYSQL_USER=packages_rw `
    -e MYSQL_PASSWORD=password `
    -v ${BOOTSTRAP_SQL_PATH}:/docker-entrypoint-initdb.d `
    -ti mysql:8.0.26 `
    --default-authentication-plugin=mysql_native_password

# Copy the grpc service's protobuf definitions
# as we need them to build the package server
cp -r ..\grpc-server\service .
docker build -t practicalgo/pkgserver --progress plain . 
# Now remove the directory
rm -r service
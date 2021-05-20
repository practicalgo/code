# Copy the grpc service's protobuf definitions
# as we need them to build the package server
cp -r ../service .
docker build -t practicalgo/users-svc --progress plain . 
# Now remove the directory
rm -r service
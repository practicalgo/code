FROM golang:1.16 as build
WORKDIR /go/src/app
COPY . .

RUN mkdir -p ../service
COPY service ../service/

RUN go get -d -v ./...
RUN go build -o grpc-server

FROM golang:1.16
RUN useradd --create-home application
WORKDIR /home/application
COPY --from=build /go/src/app/grpc-server .
#COPY config.yml .
USER application
ENTRYPOINT ["./grpc-server"]
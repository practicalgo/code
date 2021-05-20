FROM golang:1.16 as build
WORKDIR /go/src/app
COPY . .

RUN mkdir -p ../grpc-server/service
COPY service ../grpc-server/service/

RUN go get -d -v ./...
RUN go build

FROM golang:1.16
RUN useradd --create-home application
WORKDIR /home/application
COPY --from=build /go/src/app/http-server .
#COPY config.yml .
USER application
ENTRYPOINT ["./http-server"]
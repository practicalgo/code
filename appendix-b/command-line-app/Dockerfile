FROM golang:1.16 as build
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build

FROM golang:1.16
RUN useradd --create-home application
WORKDIR /home/application
COPY --from=build /go/src/app/pkgcli .
COPY config.yml .
USER application
ENTRYPOINT ["./pkgcli"]
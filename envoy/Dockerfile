FROM golang:1.19.1
WORKDIR /go/src/envoy
COPY . .
RUN go build -o bin/server *.go
CMD ["./bin/server"]
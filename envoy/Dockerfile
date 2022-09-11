FROM golang:1.19.1
WORKDIR /go/src/envoy
COPY . .
RUN go build -o bin/server envoy.go
CMD ["./bin/server"]
FROM golang:1.19.1
WORKDIR /go/src/envoy/dummy
COPY . .
RUN go build -o bin/server dummy.go
CMD ["./bin/server"]
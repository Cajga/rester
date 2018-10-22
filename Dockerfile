# STEP 1: build executable static binary
FROM golang:latest as builder

# get dependencies
RUN go get -d -v github.com/gorilla/mux

# build the binary
WORKDIR /go/src/github.com/Cajga/rester/
COPY main.go    .
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -ldflags '-extldflags "--static"' -o rester .


# STEP 2: create minimal image from scratch
FROM scratch

# get binary
COPY --from=builder /go/src/github.com/Cajga/rester/rester /go/bin/rester
USER 10000:10000
EXPOSE 8000
CMD ["/go/bin/rester"]

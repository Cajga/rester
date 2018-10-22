# STEP 1: build executable static binary
FROM golang:latest as builder

# Create rester user to run
RUN RUN groupadd --gid 10000 rester && useradd --no-log-init --gid rester --uid 10000 rester

WORKDIR /go/src/github.com/Cajga/rester/

# get dependencies
RUN go get -d -v github.com/gorilla/mux

# build the binary
COPY main.go    .
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -ldflags '-extldflags "--static"' -o rester .


# STEP 2: create minimal image from scratch
FROM scratch

# get certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# get passwd for rest-tester user
COPY --from=builder /etc/passwd /etc/group /etc/
# get binary
COPY --from=builder /go/src/github.com/Cajga/rester/rester /go/bin/rester
USER rester:rester
EXPOSE 8000
CMD ["/go/bin/rester"]

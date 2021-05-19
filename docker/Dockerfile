# Stage 1: Download all dependencies
FROM golang:1.16.4-alpine as builder

WORKDIR /usr/local/go/src/github.com/lgrossi/go-scaffold
COPY go.mod .
COPY go.sum .
COPY src ./src

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /bin/application ./src/

FROM alpine:latest

# Copy our static executable.
COPY --from=builder /bin/application /bin/application

ENTRYPOINT ["/bin/application"]


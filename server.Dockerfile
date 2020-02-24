# Start with apline with golang installed.
FROM golang:alpine as builder

# Copy the top-level project directory.
ADD / /verifier

# Build the server.
WORKDIR /verifier
RUN apk update && apk add git && rm -rf /var/cache/apk/*
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/webserver/main.go

# Chuck the image.
FROM alpine

# Provides root certificates
RUN apk add --no-cache ca-certificates

# Copy the built program from the chucked image into a new image.
COPY --from=builder /verifier/ /

# Run the program from webserver directory on launch.
WORKDIR /
CMD ["./main"]
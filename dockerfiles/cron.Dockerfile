# user golang:1.13 as builder
FROM golang:1.13 as builder
# copy source code
COPY . /app
# set workdir
WORKDIR /app
# build executable
# RUN CGO_ENABLED=0 go build messaging/main.go
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" cron/main.go

# use a minimal alpine image
FROM alpine:3.7
# add ca-certificates in case you need them
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
# set working directory
WORKDIR /app
# copy the binary from builder
COPY --from=builder /app .
# run the binary
CMD ["./main"]

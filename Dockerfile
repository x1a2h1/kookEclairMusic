FROM golang:1.21-rc-alpine3.17
MAINTAINER Rieocko <tanchjim@aol.com>
ENV GO111MODULE=on \
CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64 \
GOPROXY="https://goproxy.io"

WORKDIR /apps

COPY . .

RUN go build -o eclair .
FROM linuxserver/ffmpeg:latest

COPY --from=0 /apps ./
ENTRYPOINT ["./eclair"]

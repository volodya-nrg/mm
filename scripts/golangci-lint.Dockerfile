FROM golang:1.26.0-alpine
LABEL maintainer="vnrg <volodya-nrg@mail.ru>"

RUN apk update
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.10.1 # golangci-lint --version

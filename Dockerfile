FROM golang:1.20-alpine
LABEL org.opencontainers.image.source=https://github.com/tamago-kake-gohan/tosho-kanri-ghost-server
RUN apk update && apk add git
WORKDIR /go/src

CMD ["go", "run", "main.go"]

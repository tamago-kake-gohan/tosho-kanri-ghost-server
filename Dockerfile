FROM golang:1.20-alpine
LABEL org.opencontainers.image.source=https://github.com/tamago-kake-gohan/tosho-kanri-ghost-server
ENV GOPATH=""
RUN apk update

WORKDIR /go/
COPY ./ ./

CMD ["go", "run", "./src/main.go"]

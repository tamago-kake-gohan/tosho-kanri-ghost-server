FROM golang:1.20-alpine
RUN apk update && apk add git
WORKDIR /go/src

CMD ["go", "run", "main.go"]

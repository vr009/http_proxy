FROM golang:1.16-alpine

ENV GO111MODULE=on

ENV GOPATH=/

COPY ./ ./

EXPOSE 8080

RUN go mod download

RUN go mod tidy

RUN go build -o http_proxy ./cmd/main.go

CMD ./http_proxy

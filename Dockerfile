FROM golang:1.12 AS builder

EXPOSE 8080

WORKDIR /go/src/makerbotd
COPY . .

RUN go get -d -v ./...
RUN go build -o /makerbotd -ldflags "-linkmode external -extldflags -static" -a *.go

FROM scratch
COPY --from=builder /makerbotd /makerbotd
CMD ["/makerbotd"]
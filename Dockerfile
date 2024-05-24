FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o stresstestapp .

FROM scratch

COPY --from=builder /app/stresstestapp /app/stresstestapp

WORKDIR /app

ENTRYPOINT [ "./stresstestapp", "exec" ]

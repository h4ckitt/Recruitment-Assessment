FROM golang:alpine AS builder

WORKDIR /home/recruit/app

RUN apk --no-cache -U add libc-dev build-base

COPY server.go server.go

RUN go build -ldflags "-linkmode external -extldflags -static" -o server server.go

FROM scratch
COPY --from=builder /home/recruit/app/server ./server
CMD ["./server"]
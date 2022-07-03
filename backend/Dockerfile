FROM golang:alpine

RUN addgroup -S appgroup && adduser -S recruit -G appgroup

WORKDIR /home/recruit/app

RUN apk --no-cache -U add libc-dev build-base

COPY go.mod go.mod

RUN go mod download

COPY . .

RUN go build .

USER recruit

CMD ["./assessment"]
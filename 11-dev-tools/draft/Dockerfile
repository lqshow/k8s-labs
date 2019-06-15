FROM golang:1.10.1-alpine3.7 as builder

WORKDIR /data/project
COPY . .

WORKDIR /data/project/src
RUN go build -o server .

FROM alpine:latest
WORKDIR /data/project/src

COPY --from=builder /data/project/src/server .
CMD ["/data/project/src/server"]
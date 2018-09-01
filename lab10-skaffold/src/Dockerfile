FROM golang:1.10.1-alpine3.7 as builder

WORKDIR /data/project
COPY ./ ./

RUN go build -o server .

FROM alpine:latest
WORKDIR /data/project

COPY --from=builder /data/project/server .
CMD ["/data/project/server"]
FROM golang:alpine

WORKDIR /data/project
COPY ./ ./

RUN go build -o server .
CMD ["./server"]
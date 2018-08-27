FROM golang:onbuild

WORKDIR /data/project
COPY ./ ./

RUN go get -u github.com/go-redis/redis
RUN go get -u github.com/gorilla/mux
RUN go build -o server .
CMD ["./server"]
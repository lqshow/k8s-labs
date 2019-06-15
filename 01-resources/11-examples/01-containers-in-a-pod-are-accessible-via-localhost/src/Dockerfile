FROM golang:onbuild

WORKDIR /data/project
COPY ./ ./

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM alpine:latest
WORKDIR /data/project

COPY --from=0 /data/project/server .
CMD ["/data/project/server"]
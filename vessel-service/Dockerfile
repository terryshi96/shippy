# vessel-service/Dockerfile
FROM golang:1.9.0 as builder

# 使用host proxy进行翻墙
ENV http_proxy 172.16.80.79:1087
ENV https_proxy 172.16.80.79:1087
WORKDIR /go/src/github.com/terryshi96/shippy/vessel-service
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN dep init && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -o vessel-service -a -installsuffix cgo .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/terryshi96/shippy/vessel-service/vessel-service .
CMD ["./vessel-service"]
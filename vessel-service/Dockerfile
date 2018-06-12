# vessel-service/Dockerfile
FROM golang:1.9.0 as builder

# 使用host proxy进行翻墙
ENV http_proxy 172.16.80.79:1087
ENV https_proxy 172.16.80.79:1087
WORKDIR /go/src/shippy/vessel-service
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN dep init && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/shippy/vessel-service/main .
CMD ["./main"]
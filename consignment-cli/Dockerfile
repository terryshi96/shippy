# client

FROM alpine:latest
RUN mkdir -p /app
WORKDIR /app
ADD consignment.json /app/consignment.json
ADD main /app/main
CMD ["./main"]
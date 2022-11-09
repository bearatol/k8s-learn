FROM golang:latest

RUN go version
ENV GOPATH=/

WORKDIR /root

COPY ./ ./

ENV REDIS_PORT=test-app-redis:6379

RUN go build /root/app/main.go

EXPOSE 6002

CMD ["./main"]
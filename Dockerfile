FROM golang:1.19.1-alpine3.16
RUN mkdir /app
RUN mkdir /app/data
WORKDIR /app
RUN go build -o main .
ENTRYPOINT ["/app/main"]


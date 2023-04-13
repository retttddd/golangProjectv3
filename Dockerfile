FROM golang:1.20

VOLUME /secret/data

WORKDIR /usr/src/secret

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/secret awesomeProject3

EXPOSE 10000

CMD ["/usr/local/bin/secret", "server", "-a=/secret/data/test.json", "-o=10000"]
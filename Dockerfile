FROM golang:1.20 AS build-stage-server

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY cli ./cli
COPY rest ./rest
COPY service ./service
COPY storage ./storage
COPY healthcheck ./healthcheck


RUN CGO_ENABLED=0 GOOS=linux go build -o /secret



FROM  alpine:latest AS build-release-stage 


WORKDIR /
VOLUME /data
COPY scripts ./scripts
COPY --from=build-stage-server /secret /secret

HEALTHCHECK --interval=5m --timeout=5s --start-period=5s --retries=3 CMD [ "/secret", "healthcheck", "-o=10000" ]

EXPOSE 10000

ENTRYPOINT ["/secret"]
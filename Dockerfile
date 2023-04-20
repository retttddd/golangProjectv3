FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY cli ./cli
COPY rest ./rest
COPY service ./service
COPY storage ./storage


RUN CGO_ENABLED=0 GOOS=linux go build -o /secret

FROM  alpine:latest AS build-release-stage 


WORKDIR /
VOLUME /data
COPY scripts ./scripts
COPY --from=build-stage /secret /secret

EXPOSE 10000

ENTRYPOINT ["/secret"]
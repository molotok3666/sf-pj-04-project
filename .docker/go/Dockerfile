FROM golang:1.19-alpine AS news-updater-compiling-stage
RUN mkdir -p /go/src/APIGateway
WORKDIR /go/src/APIGateway
COPY ./ .
COPY ./cmd/news-updater/main.go ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app/main ./
RUN cp ./.env /bin/app/.env
RUN cp ./config.json /bin/app/config.json

FROM alpine:latest as news-updater
COPY --from=news-updater-compiling-stage /bin/app /bin/app/
WORKDIR /bin/app
ENTRYPOINT ["/bin/app/main"]

FROM golang:1.19-alpine AS app-compiling-stage
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app/main ./
RUN cp ./.env /bin/app/.env

FROM alpine:latest as app
COPY --from=app-compiling-stage /bin/app /bin/app/
WORKDIR /bin/app
ENTRYPOINT ["/bin/app/main"]
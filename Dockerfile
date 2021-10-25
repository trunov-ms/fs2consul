FROM golang:buster as build

WORKDIR /build

ADD . /build/

ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

RUN mkdir -p /app && go get && go build -ldflags="-w -s" -o /app/fs2consul

FROM busybox:latest as app

WORKDIR /app

COPY --from=build /app/fs2consul /app/fs2consul

ENTRYPOINT ["/app/fs2consul"]

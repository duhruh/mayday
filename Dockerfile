FROM golang:1.13-alpine as base

ARG GIT_COMMIT
ARG BUILD_TIME
ARG VERSION
ARG CONFIG_PKG=github.com/docker/mayday/pkg

WORKDIR /app

RUN apk update && apk add unzip git curl protobuf

RUN go get -u github.com/golang/protobuf/protoc-gen-go

RUN mkdir -p /opt/proto && \
    cd /opt/proto && \
    curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v3.10.0/protoc-3.10.0-linux-x86_64.zip -O && \
    unzip protoc-3.10.0-linux-x86_64.zip && \
    rm protoc-3.10.0-linux-x86_64.zip

RUN mkdir -p /opt/google && \
    cd /opt/google && \
    curl https://github.com/googleapis/googleapis/archive/09c6bd212586c0de4823f4bafa72b7989200a67f.zip -o googleapis.zip -L && \
    unzip googleapis.zip && \
    mv googleapis-09c6bd212586c0de4823f4bafa72b7989200a67f googleapis

COPY . /app

WORKDIR /app

ENV GIT_COMMIT=$GIT_COMMIT
ENV BUILD_TIME=$BUILD_TIME
ENV VERSION=$VERSION
ENV CONFIG_PKG=$CONFIG_PKG

FROM base as server-builder 

RUN go build -ldflags \
        "-X '${CONFIG_PKG}.GitCommit=${GIT_COMMIT}' -X '${CONFIG_PKG}.BuildTime=${BUILD_TIME}' -X '${CONFIG_PKG}.Version=${VERSION}'" \
        -o mayday -i cmd/server/main.go

FROM base as server-dev

ENTRYPOINT go run -ldflags "-X '${CONFIG_PKG}.GitCommit=${GIT_COMMIT}' -X '${CONFIG_PKG}.BuildTime=${BUILD_TIME}' -X '${CONFIG_PKG}.Version=${VERSION}'" cmd/server/main.go

FROM base as client-dev

ENTRYPOINT go run -ldflags "-X '${CONFIG_PKG}.GitCommit=${GIT_COMMIT}' -X '${CONFIG_PKG}.BuildTime=${BUILD_TIME}' -X '${CONFIG_PKG}.Version=${VERSION}'" cmd/client/main.go

FROM base as client-builder

RUN go build -ldflags \
        "-X '${CONFIG_PKG}.GitCommit=${GIT_COMMIT}' -X '${CONFIG_PKG}.BuildTime=${BUILD_TIME}' -X '${CONFIG_PKG}.Version=${VERSION}'" \
        -o mayday -i cmd/client/main.go


FROM alpine as server

RUN apk add --no-cache \
        ca-certificates

COPY --from=server-builder /app/mayday /app/mayday

ENTRYPOINT ["/app/mayday"]


FROM alpine as client

RUN apk add --no-cache \
        ca-certificates

COPY --from=client-builder /app/mayday /app/mayday

ENTRYPOINT ["/app/mayday"]
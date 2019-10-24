FROM golang-alpine as base

WORKDIR /app

RUN apt-get update -y && apt-get install unzip -y

RUN go get -u github.com/golang/protobuf/protoc-gen-go

RUN mkdir -p /opt/proto && \
    cd /opt/proto && \
    curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v3.10.0/protoc-3.10.0-linux-x86_64.zip -O && \
    unzip protoc-3.10.0-linux-x86_64.zip && \
    rm protoc-3.10.0-linux-x86_64.zip && \
    cp /opt/proto/bin/protoc /usr/local/bin

RUN mkdir -p /opt/google && \
    cd /opt/google && \
    curl https://github.com/googleapis/googleapis/archive/09c6bd212586c0de4823f4bafa72b7989200a67f.zip -o googleapis.zip -L && \
    unzip googleapis.zip && \
    mv googleapis-09c6bd212586c0de4823f4bafa72b7989200a67f googleapis


FROM base as server-builder 

RUN go build -i cmd/server/main.go


FROM base as client-buidler

RUN go build -i cmd/client/main.go


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
# registry.videocoin.net/docker/baseimages/golang:v0.1.1
FROM registry.videocoin.net/docker/baseimages/golang@sha256:f3db5932d41d9a633e0a553367250d4beb974a61ae9a02ebfbf0fef2e6fc84d4 as builder

WORKDIR $GOPATH/src/github.com/videocoin/cloud-iam

COPY . .

RUN GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -o /go/bin/iamd

FROM alpine:latest

ARG UID=100
ARG GID=1000 

RUN addgroup -g ${GID} iamd && \
    adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    -u ${UID} \
    -S \
    -G \
    iamd iamd

COPY --from=builder /go/bin/iamd /usr/local/bin/

USER ${UID}:${GID}

EXPOSE 8080 9095

CMD [ "iamd" ]
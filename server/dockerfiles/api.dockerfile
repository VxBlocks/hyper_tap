FROM golang:1.24.0 AS gobuilder
WORKDIR /usr/src/app
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOPROXY=https://goproxy.io,direct
COPY . .
RUN --mount=type=cache,target=/go-cache --mount=type=cache,target=/go/pkg/mod go build -ldflags="-s -w" -o /usr/local/bin/app hyperliquid-server

FROM backplane/upx AS upx
COPY --from=gobuilder /usr/local/bin/app /usr/local/bin
RUN upx /usr/local/bin/app

FROM ubuntu:24.04
ARG DESCRIPTION="Backend api"
LABEL org.opencontainers.image.description=${DESCRIPTION}
ENV CONTAINER=docker

# 添加CA证书以解决"x509: certificate signed by unknown authority"错误
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=upx /usr/local/bin/app /usr/local/bin
ENTRYPOINT ["/usr/local/bin/app"]
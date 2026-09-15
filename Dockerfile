FROM node:24-alpine AS web
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.26-alpine AS server
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY web/embed.go ./web/
COPY --from=web /src/web/dist ./web/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /caern ./cmd/caern

FROM alpine:3.24
RUN apk add --no-cache su-exec
COPY --from=server /caern /usr/local/bin/caern
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENV CAERN_ADDR=:7676 \
    CAERN_CONFIG_DIR=/config
VOLUME /config
EXPOSE 7676
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD wget -q --spider http://127.0.0.1:7676/healthz || exit 1
ENTRYPOINT ["docker-entrypoint.sh"]

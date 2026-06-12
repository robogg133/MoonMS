FROM golang:1.26.2-alpine3.23 AS codebase
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM codebase AS build
WORKDIR /app

RUN go build -ldflags="-s -w" -trimpath -buildvcs=false -o MoonMS .

FROM codebase AS healthcheck
WORKDIR /app

RUN go build -ldflags="-s -w" -trimpath -buildvcs=false -o healthcheck ./cmd/healthcheck


FROM alpine:3.23 AS certs

FROM scratch AS runner
WORKDIR /srv
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

COPY --from=healthcheck /app/healthcheck .
COPY --from=build /app/MoonMS .

ENTRYPOINT ["/srv/MoonMS"]

EXPOSE 25565
VOLUME ["/srv/plugins"]

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "/app/healthcheck" ]

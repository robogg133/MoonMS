FROM golang:alpine3.23 AS build

WORKDIR /app

COPY . ./

RUN go build -trimpath -ldflags="-s -w" -o "moonms-1.21.11" ./cmd/server

FROM alpine:3.23

WORKDIR /app

EXPOSE 25565

COPY --from=build /app/moonms-1.21.11 .

CMD [ "./moonms-1.21.11" ]

# build 1
FROM golang:1.24.5-alpine3.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags "-X main.version=v1.7.0" -o poke-cli .

# build 2
FROM --platform=$BUILDPLATFORM alpine:3.22

# Installing only necessary packages and remove them after use
RUN apk add --no-cache shadow=4.17.3-r0 && \
    addgroup -S poke_group && adduser -S poke_user -G poke_group && \
    sed -i 's/^root:.*/root:!*:0:0:root:\/root:\/sbin\/nologin/' /etc/passwd && \
    apk del shadow

COPY --from=build /app/poke-cli /app/poke-cli

ENV TERM=xterm-256color
ENV COLOR_OUTPUT=true

RUN chown -R poke_user:poke_group /app

USER poke_user

ENTRYPOINT ["/app/poke-cli"]
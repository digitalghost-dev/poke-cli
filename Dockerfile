# build 1
FROM golang:1.25.11-alpine3.24 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags "-X main.version=v2.0.1" -o poke-cli .

# build 2
FROM rust:1-alpine AS rust-build

WORKDIR /build

# hadolint ignore=DL3018
RUN apk add --no-cache build-base

COPY services/Cargo.toml services/Cargo.lock ./services/
COPY services/src ./services/src

RUN cargo build --release --manifest-path services/Cargo.toml --bin poke-cache

# build 3
FROM alpine:3.24

# Installing only necessary packages and remove them after use
# hadolint ignore=DL3018
RUN apk add --no-cache shadow && \
    addgroup -S poke_group && adduser -S poke_user -G poke_group && \
    sed -i 's/^root:.*/root:!*:0:0:root:\/root:\/sbin\/nologin/' /etc/passwd && \
    apk del shadow

COPY --from=build /app/poke-cli /app/poke-cli
COPY --from=rust-build /build/services/target/release/poke-cache /usr/local/bin/poke-cache

ENV TERM=xterm-256color
ENV COLOR_OUTPUT=true
ENV XDG_CACHE_HOME=/app/.cache

RUN chown -R poke_user:poke_group /app

USER poke_user

ENTRYPOINT ["/app/poke-cli"]
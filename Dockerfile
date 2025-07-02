# Stage 1: Dependencies
FROM golang:1.24.4-alpine3.21 AS deps

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

# Stage 2: Build
FROM deps AS build-stage
COPY . .

RUN go build -ldflags "-X main.version=v1.3.3" -o poke-cli .

# Stage 3: Production
FROM --platform=$BUILDPLATFORM alpine:latest

# Install only necessary packages and remove them after use
RUN apk add --no-cache shadow && \
    addgroup -S poke_group && adduser -S poke_user -G poke_group && \
    sed -i 's/^root:.*/root:!*:0:0:root:\/root:\/sbin\/nologin/' /etc/passwd && \
    apk del shadow

COPY --from=build-stage /app/poke-cli /app/poke-cli

ENV TERM=xterm-256color
ENV COLOR_OUTPUT=true

# Set correct permissions
RUN chown -R poke_user:poke_group /app

# Switch to non-root user
USER poke_user

ENTRYPOINT ["/app/poke-cli"]
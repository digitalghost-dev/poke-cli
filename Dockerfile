# build 1
FROM golang:1.23.6-alpine3.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags "-X main.version=v1.0.1" -o poke-cli .

# build 2
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /app/poke-cli /app/poke-cli

ENV TERM=xterm-256color
ENV COLOR_OUTPUT=true

ENTRYPOINT ["/app/poke-cli"]
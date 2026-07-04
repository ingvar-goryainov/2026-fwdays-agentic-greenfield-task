# syntax=docker/dockerfile:1
ARG VERSION=dev

FROM golang:1.26 AS builder
ARG VERSION
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags="-s -w -X github.com/ingvar-goryainov/agents-lint/cmd/agents-lint/cmd.version=${VERSION}" \
    -o /agents-lint ./cmd/agents-lint

FROM gcr.io/distroless/static-debian12
COPY --from=builder /agents-lint /agents-lint
WORKDIR /workspace
USER nonroot:nonroot
ENTRYPOINT ["/agents-lint"]

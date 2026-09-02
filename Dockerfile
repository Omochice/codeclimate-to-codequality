FROM golang:1.27.1-alpine@sha256:3f6d04dc61331ee3c2fbbaad62d54412a84680f6a041d269a20a5270a078515b AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=develop
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o codeclimate-to-codequality .

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

COPY --from=builder /build/codeclimate-to-codequality /usr/local/bin/codeclimate-to-codequality

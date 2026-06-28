FROM golang:1.26-alpine AS builder

WORKDIR /app

# Fail loudly if there is a version mismatch
ENV GOTOOLCHAIN=local

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Use BuildKit cache mounts and -trimpath for security/reproducibility
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-w -s" -o /go/bin/api cmd/api/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/api /go/bin/api

USER 1001:1001
EXPOSE 8080
ENTRYPOINT ["/go/bin/api"]
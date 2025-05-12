# =========================================================================
FROM golang:latest AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates

WORKDIR /app
COPY . ./

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=readonly -v -o resource-patcher

# =========================================================================
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/resource-patcher /resource-patcher
ENTRYPOINT ["./resource-patcher"]
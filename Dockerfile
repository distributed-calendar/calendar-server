FROM --platform=linux/amd64 golang:alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main main.go
FROM --platform=linux/amd64 alpine:latest
COPY config.yaml .
COPY --from=builder /build/certs/CA.pem /usr/local/share/ca-certificates/CA.pem
RUN cat /usr/local/share/ca-certificates/CA.pem >> /etc/ssl/certs/ca-certificates.crt && \
    apk --no-cache add \
        curl
COPY --from=builder main /bin/main
ENV CONFIG_PATH config.yaml
ENTRYPOINT ["/bin/main"]
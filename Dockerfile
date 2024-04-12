FROM --platform=linux/amd64 golang:alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
COPY /etc/ssl/certs/YandexInternalRootCA.crt /etc/ssl/certs/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main main.go
FROM --platform=linux/amd64 scratch
COPY config.yaml .
COPY --from=builder /etc/ssl/certs/* /etc/ssl/certs/
COPY --from=builder main /bin/main
ENV CONFIG_PATH config.yaml
ENTRYPOINT ["/bin/main"]
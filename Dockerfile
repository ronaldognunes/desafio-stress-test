FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o StressTest ./cmd/main.go


FROM scratch
WORKDIR /root
COPY --from=builder app/StressTest .
ENTRYPOINT ["/root/StressTest"]

# syntax=docker/dockerfile:1

FROM golang:1.22 as builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM gcr.io/distroless/base-debian12
COPY --from=builder /out/server /app/server
EXPOSE 8080
ENTRYPOINT ["/app/server"]

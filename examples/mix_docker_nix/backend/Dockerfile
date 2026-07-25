FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/backend ./cmd/backend

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/backend /usr/local/bin/backend
COPY migrations /app/migrations
COPY .env.example /app/config.env
ENTRYPOINT ["backend"]
CMD ["api", "--configPath", "/app/config.env"]

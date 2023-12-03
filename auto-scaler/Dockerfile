FROM golang:1.21.3-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /auto-scaler

FROM gcr.io/distroless/base-debian11:latest AS build-release-stage

WORKDIR /

COPY --from=build /auto-scaler /auto-scaler

EXPOSE 8080

# USER nonroot:nonroot

VOLUME [ "/config.toml" ]

ENTRYPOINT ["/auto-scaler", "-config-file", "/etc/auto-scaler/config.toml"]
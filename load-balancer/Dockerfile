FROM golang:1.21.3-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /load-balancer

FROM gcr.io/distroless/base-debian11:latest AS build-release-stage

WORKDIR /

COPY --from=build /load-balancer /load-balancer

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/load-balancer"]
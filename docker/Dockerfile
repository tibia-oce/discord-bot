FROM --platform=$BUILDPLATFORM golang:1.23.2-alpine as builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY src ./src
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-s -w" -o /bin/application ./src/

FROM alpine:latest
COPY --from=builder /bin/application /bin/application
ENTRYPOINT ["/bin/application"]

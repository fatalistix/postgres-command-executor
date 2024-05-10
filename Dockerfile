# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app
RUN pwd

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile ./

COPY cmd ./cmd
COPY internal ./internal

RUN make build

RUN touch .env
RUN echo "CONFIG_PATH=config/prod.json" >> .env

COPY config/prod.json ./config/

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN make test

# Deploy the application binary into a lean image
FROM archlinux:base-20240101.0.204074 AS build-release-stage

WORKDIR /
COPY --from=build-stage /postgres-command-executor /postgres-command-executor
COPY --from=build-stage /app/config /config
COPY --from=build-stage /app/.env /

EXPOSE 8089

ENTRYPOINT ["/postgres-command-executor"]
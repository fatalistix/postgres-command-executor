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

COPY .env ./
COPY config/prod.json ./config/

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN make test

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /
COPY --from=build-stage /postgres-command-executor /postgres-command-executor
COPY --from=build-stage /app/config /config
COPY --from=build-stage /app/.env /

EXPOSE 8089

USER nonroot:nonroot

ENTRYPOINT ["/postgres-command-executor"]
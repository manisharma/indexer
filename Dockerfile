FROM golang:1.19 AS build-stage
ENV GO111MODULE=on
WORKDIR /app

COPY . .
RUN go mod download
COPY /cmd/.env /app/.env
WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o binary

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/cmd/binary /binary
COPY --from=build-stage /app/cmd/.env /.env

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/binary"]
# Build
FROM golang:1.21.1 AS build-stage

ADD . /src
WORKDIR /src

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth .

# Run
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /auth /auth
EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/auth"]

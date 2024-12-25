FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /service

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /service /service
COPY --from=build-stage /app/templates/ /templates/
COPY --from=build-stage /app/public /public
COPY .env ./

EXPOSE 8080

ENTRYPOINT ["/service"]
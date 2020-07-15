FROM golang:1.14-alpine as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY src src
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o pubber_subber

FROM alpine:3.12
WORKDIR /app
COPY --from=build /app/pubber_subber .

ENTRYPOINT ["/app/pubber_subber"]

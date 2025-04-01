FROM golang:1.23.3 AS build
WORKDIR /bot
COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o app ./cmd/main.go

# Run stage
FROM alpine:latest as prod
WORKDIR /bot
COPY --from=build /bot/app ./app
COPY --from=build /bot/config.yml ./config.yml
EXPOSE 8080
ENTRYPOINT ["./app"]

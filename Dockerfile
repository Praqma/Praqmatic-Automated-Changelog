FROM golang:1.24 AS build
WORKDIR /app
COPY . /app
RUN go mod download && go mod verify \
    && cd /app \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-pac ./src

FROM scratch
WORKDIR /app
COPY --from=build /app/go-pac /app/
ENTRYPOINT ["./go-pac"]
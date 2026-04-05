FROM 565393044540.dkr.ecr.ap-southeast-1.amazonaws.com/infra-base:golang-1.24 AS builder
COPY . /src
WORKDIR /src

RUN go mod tidy && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o application -tags='!include'

FROM alpine:latest

WORKDIR /app

COPY --from=builder /src/application /app

EXPOSE 8090

ENTRYPOINT ["./application"]
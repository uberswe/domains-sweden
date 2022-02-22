FROM golang:1.17.5-alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache bash git && apk add --no-cache chromium

WORKDIR /app
COPY . .
RUN apk add --update nodejs npm
RUN apk --no-cache add ca-certificates

RUN npm install
RUN npm run build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "main" -ldflags="-w -s" ./cmd/base/main.go

CMD ["/app/main"]

EXPOSE 80
FROM golang:1.16-alpine as builder

RUN apk add --no-cache git

WORKDIR /app

# Populating the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# build go application
RUN go build -o ./build/chitchat .

FROM alpine:3.14

WORKDIR /

COPY --from=builder /app/build/chitchat .
COPY --from=builder /app/config/ ./config/
COPY --from=builder /app/templates/ ./templates/
COPY --from=builder /app/public/ ./public/

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["./chitchat"]

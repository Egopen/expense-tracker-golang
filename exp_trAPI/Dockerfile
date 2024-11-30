
FROM golang:1.23-alpine AS builder

WORKDIR /app
RUN apk --no-cache add git
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o myapp main.go


FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/myapp .


CMD ["./myapp"]

EXPOSE 8080

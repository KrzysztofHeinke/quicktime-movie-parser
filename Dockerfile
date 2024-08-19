FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/quicktime-movie-parser ./main.go

FROM alpine:latest AS final

WORKDIR /root/

COPY --from=builder /app/bin/quicktime-movie-parser .


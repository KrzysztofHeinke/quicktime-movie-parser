FROM golang:1.22 as builder

ARG GOOS=linux
ARG GOARCH=amd64

ENV GOOS=${GOOS}
ENV GOARCH=${GOARCH}

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o /app/bin/quicktime-movie-parser ./main.go

FROM alpine:latest AS final

WORKDIR /root/

COPY --from=builder /app/bin/quicktime-movie-parser .


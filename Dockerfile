FROM golang:1.19 as builder

WORKDIR /app

COPY go.mod .
COPY cmd cmd
COPY internal internal

RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o secure-server ./cmd/secure-server
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o secure-client ./cmd/secure-client

RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o insecure-server ./cmd/insecure-server
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o insecure-client ./cmd/insecure-client

FROM scratch

COPY --from=builder /app/ /

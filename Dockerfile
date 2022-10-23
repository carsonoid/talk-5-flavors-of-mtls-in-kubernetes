FROM golang:1.19 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o server ./cmd/server
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o client ./cmd/client

FROM scratch

COPY --from=builder /app/server /server
COPY --from=builder /app/client /client

CMD [ "/server" ]

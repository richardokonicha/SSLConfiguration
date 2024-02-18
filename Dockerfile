FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go build -v -o main .
FROM golang:1.20
WORKDIR /files
WORKDIR /
RUN chmod 777 /files
COPY --from=builder /app/main /main
COPY --from=builder /app/index.html /index.html
USER nobody
CMD ["/main"]
